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
values ($1, $2, $3, $4, $5, $6, $7, $8 returning id);`
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
values ($1, $2, $3, $4, $5, $6, $7 returning id);`
	var id uuid.UUID
	err := r.db.QueryRow(query, transaction.BudgetID, userID, transaction.CategoryID, transaction.Amount, transaction.Intent, transaction.Direction, transaction.Comment).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to insert transaction: %w", err)
	}
	return id, nil

}
func (r *TransactionPostgres) UpdateTransactionWithGoal(transaction models.UpdateTransaction, old_amount decimal.Decimal) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := `update budget.transactions 
set category_id = $1, 
goal_id = $2,
amount = $3,
date_update = current_timestamp,
intent = $4,
direction = $5,
comment = $6
where transaction_id = $7 ;`
	_, err = tx.Exec(query, transaction.CategoryID, transaction.GoalID, transaction.Amount, transaction.Intent, transaction.Direction, transaction.Comment, transaction.TransactionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	query2 := `UPDATE budget.accumulation_goals 
			SET current_saved = current_saved - $1 + $2
			WHERE uuid = $3;`
	_, err = tx.Exec(query2, old_amount, transaction.Amount, transaction.GoalID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update goal balance: %w", err)
	}

	return tx.Commit()
}
func (r *TransactionPostgres) UpdateTransactionWithoutGoal(transaction models.UpdateTransaction) error {
	query := `update budget.transactions 
set category_id = NULL, 
goal_id = $1,
amount = $2,
date_update = current_timestamp,
intent = $3,
direction = $4,
comment = $5
where transaction_id = $6 ;`
	_, err := r.db.Exec(query, transaction.GoalID, transaction.Amount, transaction.Intent, transaction.Direction, transaction.Comment, transaction.TransactionID)
	if err != nil {
		return fmt.Errorf("failed to update transaction: %w", err)
	}
	return nil
}

func (r *TransactionPostgres) DeleteTransactionWithoutGoal(transactionID uuid.UUID) error {
	query := `delete from budget.transactions where transaction_id = $1`
	_, err := r.db.Exec(query, transactionID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	return nil
}

func (r *TransactionPostgres) DeleteTransactionWithGoal(transactionID uuid.UUID, amount decimal.Decimal, goalID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := `delete from budget.transactions where transaction_id = $1`
	_, err = tx.Exec(query, transactionID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to delete transaction: %w", err)
	}
	query2 := `UPDATE budget.accumulation_goals 
			SET current_saved = current_saved - $1
			WHERE uuid = $2;`
	_, err = tx.Exec(query2, amount, goalID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update goal balance: %w", err)
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

func (r *TransactionPostgres) GetTransactionAmountByID(transactionID uuid.UUID) (decimal.Decimal, error) {
	var amount decimal.Decimal
	query := `select amount from budget.transactions where transaction_id = $1`
	err := r.db.QueryRow(query, transactionID).Scan(&amount)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to get transaction amount: %w", err)
	}
	return amount, nil
}
