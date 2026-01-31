package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type NewGoal struct {
	BudgetID     uuid.UUID       `json:"budget_id"`
	UserID       uuid.UUID       `*`
	TargetName   string          `json:"target_name"`
	TargetAmount decimal.Decimal `json:"target_amount"`
	TargetDate   time.Time       `json:"target_date"`
	CurrencyID   int             `json:"currency_id"`
}

type UpdateGoal struct {
	ID           uuid.UUID       `json:"id"`
	TargetName   string          `json:"target_name"`
	TargetAmount decimal.Decimal `json:"target_amount"`
	TargetDate   time.Time       `json:"target_date"`
	CurrencyID   int             `json:"currency_id"`
}
type Goal struct {
	ID                  uuid.UUID       `json:"id"`
	TargetName          string          `json:"target_name"`
	TargetAmount        decimal.Decimal `json:"target_amount"`
	TargetDate          time.Time       `json:"target_date"`
	CurrentSaved        decimal.Decimal `json:"current_saved"`
	RequiredMonthlySave decimal.Decimal `json:"required_monthly_save"`
	CurrencyCode        string          `json:"currency_code"`
	CurrencyName        string          `json:"currency_name"`
}

type GoalID struct {
	ID uuid.UUID `json:"id"`
}
