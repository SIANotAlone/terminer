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
	query := `insert into budget.transactions(budget_id, user_id, category_id, goal_id, amount, intent, direction, comment)
values ($1, $2, $3, $4, $5, $6, $7, $8 )returning transaction_id;`
	var id uuid.UUID
	err = tx.QueryRow(query, transaction.BudgetID, userID, transaction.CategoryID, transaction.GoalID, transaction.Amount, transaction.Intent, transaction.Direction, transaction.Comment).Scan(&id)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to insert transaction: %w", err)
	}
	query2 := `UPDATE budget.accumulation_goals 
            SET current_saved = current_saved + $1 
            WHERE uuid = $2;`
	_, err = tx.Exec(query2, transaction.Amount, transaction.GoalID)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, fmt.Errorf("failed to update goal balance: %w", err)
	}

	return id, tx.Commit()

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

func (r *TransactionPostgres) UpdateTransaction(transaction models.UpdateTransaction, oldAmount decimal.Decimal) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    // Обновляем всё: и категорию, и сумму, и прочее.
    // Через RETURNING забираем goal_id, который УЖЕ лежит в базе.
    query := `
        UPDATE budget.transactions 
        SET category_id = $1, 
            amount = $2,
            intent = $3,
            direction = $4,
            comment = $5,
            date_update = current_timestamp
        WHERE transaction_id = $6
        RETURNING goal_id;`

    var goalID uuid.NullUUID
    err = tx.QueryRow(query, 
        transaction.CategoryID, 
        transaction.Amount, 
        transaction.Intent, 
        transaction.Direction, 
        transaction.Comment, 
        transaction.TransactionID,
    ).Scan(&goalID)

    if err != nil {
        tx.Rollback()
        return fmt.Errorf("update failed: %w", err)
    }

    // Если у транзакции БЫЛА цель (неважно, какая категория), обновляем её баланс
    if goalID.Valid {
        query2 := `UPDATE budget.accumulation_goals 
                   SET current_saved = current_saved - $1 + $2
                   WHERE uuid = $3;`
        _, err = tx.Exec(query2, oldAmount, transaction.Amount, goalID.UUID)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("goal balance update failed: %w", err)
        }
    }

    return tx.Commit()
}
func (r *TransactionPostgres) DeleteTransaction(transactionID uuid.UUID) error {
    tx, err := r.db.Begin()
    if err != nil {
        return err
    }

    // 1. Удаляем транзакцию и СРАЗУ возвращаем её сумму и цель, если они были.
    // Это атомарная операция: мы берем данные именно той строки, которую удаляем.
    query := `
        DELETE FROM budget.transactions 
        WHERE transaction_id = $1 
        RETURNING amount, goal_id;`

    var amount decimal.Decimal
    var goalID uuid.NullUUID

    err = tx.QueryRow(query, transactionID).Scan(&amount, &goalID)
    if err != nil {
        tx.Rollback()
        // Если транзакция не найдена, Scan вернет sql.ErrNoRows
        return fmt.Errorf("transaction not found or failed to delete: %w", err)
    }

    // 2. Если у удаленной транзакции БЫЛА цель, корректируем её баланс.
    if goalID.Valid {
        query2 := `
            UPDATE budget.accumulation_goals 
            SET current_saved = current_saved - $1
            WHERE uuid = $2;`
        
        _, err = tx.Exec(query2, amount, goalID.UUID)
        if err != nil {
            tx.Rollback()
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
