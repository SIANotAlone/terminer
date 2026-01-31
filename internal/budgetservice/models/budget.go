package models

import (
	"time"

	"github.com/google/uuid"
)

type NewBudget struct {
	Name          string    `json:"name" binding:"required" omitempty:"true"`
	Type          int       `json:"type" binding:"required" omitempty:"true"`
	Date_Start    time.Time `json:"date_start" binding:"required" omitempty:"true"`
	Date_End      time.Time `json:"date_end" binding:"required" omitempty:"true"`
	Base_Currency string    `json:"base_currency" binding:"required" omitempty:"true"`
}

type BudgetType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Budget struct {
	ID            uuid.UUID `json:"id"`
	Owner         uuid.UUID `json:"owner"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Date          time.Time `json:"date"`
	Date_Start    time.Time `json:"date_start"`
	Date_End      time.Time `json:"date_end"`
	CurrencyCode  string    `json:"currency_code"`
	Base_Currency string    `json:"base_currency"`
	Is_Archived   bool      `*`
}

type UpdateBudget struct {
	ID              uuid.UUID `json:"id" binding:"required" omitempty:"true"`
	Name            string    `json:"name" binding:"required" omitempty:"true"`
	TypeID          int       `json:"type_id" binding:"required" omitempty:"true"`
	Date_Start      time.Time `json:"date_start" binding:"required" omitempty:"true"`
	Date_End        time.Time `json:"date_end" binding:"required" omitempty:"true"`
	Base_CurrencyID int       `json:"base_currency" binding:"required" omitempty:"true"`
}

type BudgetID struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}