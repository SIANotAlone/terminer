package repository

import (
	"fmt"
	"terminer/internal/budgetservice/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type TransactionPostgres struct {
	db *sqlx.DB
}

func NewTransactionPostgres(db *sqlx.DB) *TransactionPostgres {
	return &TransactionPostgres{db: db}
}

func (r *TransactionPostgres) CreateTransactionWithGoal(userID uuid.UUID, transaction models.NewTransaction) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback()

	query := `
		INSERT INTO budget.transactions(
			budget_id, user_id, category_id, goal_id, amount, intent, direction, comment
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING transaction_id;
	`

	var id uuid.UUID
	err = tx.QueryRow(
		query,
		transaction.BudgetID,
		userID,
		transaction.CategoryID,
		transaction.GoalID,
		transaction.Amount,
		transaction.Intent,
		transaction.Direction,
		transaction.Comment,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert transaction: %w", err)
	}

	delta := transaction.Amount
	switch transaction.Direction {
	case "EXPENSE":
		delta = transaction.Amount
	case "INCOME":
		delta = transaction.Amount.Neg()
	default:
		return uuid.Nil, fmt.Errorf("unknown transaction direction: %s", transaction.Direction)
	}

	query2 := `
		UPDATE budget.accumulation_goals
		SET current_saved = current_saved + $1
		WHERE uuid = $2;
	`

	_, err = tx.Exec(query2, delta, transaction.GoalID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to update goal balance: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return uuid.Nil, err
	}

	return id, nil
}
func (r *TransactionPostgres) CreateTransactionWithoutGoal(userID uuid.UUID, transaction models.NewTransaction) (uuid.UUID, error) {
	query := `insert into budget.transactions(budget_id, user_id, category_id, amount, intent, direction, comment)
values ($1, $2, $3, $4, $5, $6, $7) returning transaction_id;`
	var id uuid.UUID
	err := r.db.QueryRow(query, transaction.BudgetID, userID, transaction.CategoryID, transaction.Amount, transaction.Intent, transaction.Direction, transaction.Comment).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert transaction: %w", err)
	}
	return id, nil

}
func goalSignedDelta(amount decimal.Decimal, direction string) (decimal.Decimal, error) {
	switch direction {
	case "EXPENSE":
		// расход -> в накоплениях это плюс
		return amount, nil
	case "INCOME":
		// доход -> деньги вернулись в бюджет, значит из накоплений вычитаем
		return amount.Neg(), nil
	default:
		return decimal.Decimal{}, fmt.Errorf("unknown direction: %s", direction)
	}
}

func (r *TransactionPostgres) UpdateTransaction(userID uuid.UUID, newTx models.UpdateTransaction) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var oldAmount decimal.Decimal
	var oldGoal uuid.NullUUID
	var oldDirection string
	var ownerID uuid.UUID

	qSelect := `
		SELECT amount, goal_id, user_id, direction
		FROM budget.transactions
		WHERE transaction_id = $1
		FOR UPDATE;
	`
	if err = tx.QueryRow(qSelect, newTx.TransactionID).Scan(&oldAmount, &oldGoal, &ownerID, &oldDirection); err != nil {
		return fmt.Errorf("select transaction for update failed: %w", err)
	}

	if ownerID != userID {
		return fmt.Errorf("forbidden")
	}

	// обновляем саму транзакцию
	var newGoal uuid.NullUUID
	qUpdate := `
		UPDATE budget.transactions
		SET category_id = $1,
			amount = $2,
			intent = $3,
			direction = $4,
			comment = $5,
			goal_id = $6,
			date_update = current_timestamp
		WHERE transaction_id = $7
		RETURNING goal_id;
	`
	if err = tx.QueryRow(
		qUpdate,
		newTx.CategoryID,
		newTx.Amount,
		newTx.Intent,
		newTx.Direction,
		newTx.Comment,
		newTx.GoalID,
		newTx.TransactionID,
	).Scan(&newGoal); err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	oldDelta, err := goalSignedDelta(oldAmount, oldDirection)
	if err != nil {
		return err
	}

	newDelta, err := goalSignedDelta(newTx.Amount, newTx.Direction)
	if err != nil {
		return err
	}

	// корректируем цели с учетом старой/новой цели и direction
	switch {
	case oldGoal.Valid && newGoal.Valid && oldGoal.UUID == newGoal.UUID:
		// та же цель: меняем только разницу между старым и новым вкладом
		q := `UPDATE budget.accumulation_goals
			  SET current_saved = current_saved - $1 + $2
			  WHERE uuid = $3;`
		_, err = tx.Exec(q, oldDelta, newDelta, newGoal.UUID)

	case oldGoal.Valid && newGoal.Valid && oldGoal.UUID != newGoal.UUID:
		// цель поменялась: убрать старый вклад из старой, добавить новый вклад в новую
		q1 := `UPDATE budget.accumulation_goals
			   SET current_saved = current_saved - $1
			   WHERE uuid = $2;`
		_, err = tx.Exec(q1, oldDelta, oldGoal.UUID)
		if err != nil {
			return fmt.Errorf("balance adjustment failed: %w", err)
		}

		q2 := `UPDATE budget.accumulation_goals
			   SET current_saved = current_saved + $1
			   WHERE uuid = $2;`
		_, err = tx.Exec(q2, newDelta, newGoal.UUID)

	case oldGoal.Valid && !newGoal.Valid:
		// цель была, теперь нет
		q := `UPDATE budget.accumulation_goals
			  SET current_saved = current_saved - $1
			  WHERE uuid = $2;`
		_, err = tx.Exec(q, oldDelta, oldGoal.UUID)

	case !oldGoal.Valid && newGoal.Valid:
		// цели не было, теперь есть
		q := `UPDATE budget.accumulation_goals
			  SET current_saved = current_saved + $1
			  WHERE uuid = $2;`
		_, err = tx.Exec(q, newDelta, newGoal.UUID)
	}

	if err != nil {
		return fmt.Errorf("balance adjustment failed: %w", err)
	}

	return tx.Commit()
}
func (r *TransactionPostgres) DeleteTransaction(transactionID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `
		DELETE FROM budget.transactions 
		WHERE transaction_id = $1 
		RETURNING amount, goal_id, direction;
	`

	var amount decimal.Decimal
	var goalID uuid.NullUUID
	var direction string

	err = tx.QueryRow(query, transactionID).Scan(&amount, &goalID, &direction)
	if err != nil {
		return fmt.Errorf("transaction not found or failed to delete: %w", err)
	}

	if goalID.Valid {
		var delta decimal.Decimal

		switch direction {
		case "EXPENSE":
			// Было +amount, при удалении убираем это из цели
			delta = amount
		case "INCOME":
			// Было -amount, при удалении возвращаем обратно
			delta = amount.Neg()
		default:
			return fmt.Errorf("unknown direction: %s", direction)
		}

		query2 := `
			UPDATE budget.accumulation_goals 
			SET current_saved = current_saved - $1
			WHERE uuid = $2;
		`
		_, err = tx.Exec(query2, delta, goalID.UUID)
		if err != nil {
			return fmt.Errorf("failed to update goal balance: %w", err)
		}
	}

	return tx.Commit()
}

func (r *TransactionPostgres) GetTransactionsByBudget(budgetID uuid.UUID) ([]models.Transaction, error) {
	query := `SELECT dc.transaction_id, u.first_name || ' '|| u.last_name as user, c.name as category, 
COALESCE(g.target_name, 'Немає') as goal, dc.amount, dc.date, dc.date_update, 
dc.intent, dc.direction, dc.comment
FROM budget.transactions dc
LEFT JOIN main.user u on dc.user_id = u.uuid 
LEFT JOIN budget.categories c on c.uuid = dc.category_id 
LEFT JOIN budget.accumulation_goals g on g.uuid = dc.goal_id 
WHERE dc.budget_id = $1;`
	rows, err := r.db.Query(query, budgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}
	defer rows.Close()
	var transactions []models.Transaction
	for rows.Next() {
		var transaction models.Transaction
		err := rows.Scan(&transaction.ID, &transaction.User, &transaction.Category, &transaction.Goal, &transaction.Amount,
			&transaction.Date, &transaction.DateUpdate, &transaction.Intent, &transaction.Direction, &transaction.Comment)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, nil

}

func (r *TransactionPostgres) HasGoal(transactionID uuid.UUID) (bool, error) {
	var goalID uuid.NullUUID
	query := `select goal_id from budget.transactions where transaction_id = $1`
	err := r.db.QueryRow(query, transactionID).Scan(&goalID)
	if err != nil {
		return false, fmt.Errorf("failed to get transaction goal: %w", err)
	}
	if goalID.Valid {
		return true, nil
	}
	return false, nil
}
func (r *TransactionPostgres) GetTrasactionWihtGoalAmount(transactionID uuid.UUID) (decimal.Decimal, uuid.UUID, error) {
	{
		var amount decimal.Decimal
		var goalID uuid.UUID
		query := `select amount, goal_id from budget.transactions where transaction_id = $1`
		err := r.db.QueryRow(query, transactionID).Scan(&amount, &goalID)
		if err != nil {
			return decimal.Decimal{}, uuid.Nil, fmt.Errorf("failed to get transaction amount and goal: %w", err)
		}
		return amount, goalID, nil
	}
}

func (r *TransactionPostgres) GetTransactionAmountByID(transactionID uuid.UUID) (decimal.Decimal, error) {
	var amount decimal.Decimal
	query := `select amount from budget.transactions where transaction_id = $1`
	err := r.db.QueryRow(query, transactionID).Scan(&amount)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to get transaction amount: %w", err)
	}
	return amount, nil
}

func (r *TransactionPostgres) GetBudgetIdByTransactionID(transactionID uuid.UUID) (uuid.UUID, error) {
	var budgetID uuid.UUID
	query := `select budget_id from budget.transactions where transaction_id = $1`
	err := r.db.QueryRow(query, transactionID).Scan(&budgetID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get budget id by transaction id: %w", err)
	}
	return budgetID, nil
}
