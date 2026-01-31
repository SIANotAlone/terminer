package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type NewTransaction struct {
	BudgetID   uuid.UUID       `json:"budget_id"`
	UserID     uuid.UUID       `*`
	CategoryID uuid.UUID       `json:"category_id"`
	GoalID     *uuid.UUID      `json:"goal_id"`
	Amount     decimal.Decimal `json:"amount"`
	Intent     string          `json:"intent"`
	Direction  string          `json:"direction"`
	Comment    string          `json:"comment"`
}

type UpdateTransaction struct {
	TransactionID uuid.UUID       `json:"transaction_id"`
	CategoryID    uuid.UUID       `json:"category_id"`
	GoalID        *uuid.UUID      `json:"goal_id"`
	Amount        decimal.Decimal `json:"amount"`
	Intent        string          `json:"intent"`
	Direction     string          `json:"direction"`
	Comment       string          `json:"comment"`
}

type Transaction struct {
	ID         uuid.UUID       `json:"id"`
	User       string          `json:"user"`
	Category   string          `json:"category"`
	Goal       *string         `json:"goal"`
	Amount     decimal.Decimal `json:"amount"`
	Date       time.Time       `json:"date"`
	DateUpdate time.Time       `json:"date_update"`
	Intent     string          `json:"intent"`
	Direction  string          `json:"direction"`
	Comment    string          `json:"comment"`
}

type TransactionID struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}

