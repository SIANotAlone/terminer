package repository

import (
	"terminer/internal/budgetservice/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type BudgetPostgres struct {
	db *sqlx.DB
}

func NewBudgetPostgres(db *sqlx.DB) *BudgetPostgres {
	return &BudgetPostgres{db: db}
}

func (r *BudgetPostgres) CreateBudget(userID uuid.UUID, budget models.NewBudget) (uuid.UUID, error) {
	query := `
INSERT INTO budget.budgets (owner_id, name, type, date, date_start, date_end, base_currency)
VALUES ($1, $2, $3, now(), $4, $5, $6)
RETURNING uuid;
`
	var budgetID uuid.UUID
	err := r.db.QueryRow(query, userID, budget.Name, budget.Type, budget.Date_Start, budget.Date_End, budget.Base_Currency).Scan(&budgetID)
	if err != nil {
		return uuid.Nil, err
	}
	return budgetID, nil
}
func (r *BudgetPostgres) UpdateBudget(userID uuid.UUID, budget models.UpdateBudget) error {
	query := `update budget.budgets 
set name= $1,type= $2, date_start=$3, date_end=$4, base_currency=$5
where uuid =$6;`
	_, err := r.db.Exec(query, budget.Name, budget.TypeID, budget.Date_Start, budget.Date_End, budget.Base_CurrencyID, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetPostgres) DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	query := `
delete from budget.budgets where uuid = $1;`
	_, err := r.db.Exec(query, budgetID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetPostgres) ArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	query := `
UPDATE budget.budgets 
SET archived = true 
WHERE uuid = $1;`
	_, err := r.db.Exec(query, budgetID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetPostgres) UnArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	query := `
UPDATE budget.budgets 
SET archived = false 
WHERE uuid = $1;`
	_, err := r.db.Exec(query, budgetID)
	if err != nil {
		return err
	}
	return nil
}

func (r *BudgetPostgres) GetAvailableBudgets(userID uuid.UUID) ([]models.Budget, error) {
	query := `
SELECT b.uuid, b.name, bt.name as type, b.date, b.date_start, b.date_end, bc.code, bc.name as currency, b.archived
FROM budget.budgets b
LEFT JOIN budget.budget_types bt on bt.id = b.type
LEFT JOIN budget.currencies bc on bc.id = b.base_currency
WHERE 
    -- 1. Пользователь является владельцем
    (b.owner_id = $1 AND b.archived = false)
    
    OR 
    
    -- 2. Бюджет расшарен пользователю через таблицу access
    (b.uuid IN (
        SELECT a.budget_id 
        FROM budget.access a 
        WHERE a.user_id = $1
    )AND b.archived=false);
`
	var budgets []models.Budget
	var budget models.Budget
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err := rows.Scan(
			&budget.ID,
			&budget.Name,
			&budget.Type,
			&budget.Date,
			&budget.Date_Start,
			&budget.Date_End,
			&budget.CurrencyCode,
			&budget.Base_Currency,
			&budget.Is_Archived,
		)
		if err != nil {
			return nil, err
		}
		budgets = append(budgets, budget)
	}
	return budgets, nil
}

func (r *BudgetPostgres) GetBudgetTypes() ([]models.BudgetType, error) {
	query := `select id, name, description from budget.budget_types;`
	var budgetTypes []models.BudgetType
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var budgetType models.BudgetType
		err := rows.Scan(&budgetType.ID, &budgetType.Name, &budgetType.Description)
		if err != nil {
			return nil, err
		}
		budgetTypes = append(budgetTypes, budgetType)
	}
	return budgetTypes, nil
}

func (r *BudgetPostgres) GetBudgetOwnerID(budgetID uuid.UUID) (uuid.UUID, error) {
	query := `select owner_id from budget.budgets where uuid = $1;`
	var ownerID uuid.UUID
	err := r.db.QueryRow(query, budgetID).Scan(&ownerID)
	if err != nil {
		return uuid.Nil, err
	}
	return ownerID, nil
}
