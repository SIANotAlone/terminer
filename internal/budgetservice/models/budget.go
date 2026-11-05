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
	Base_Currency int       `json:"base_currency" binding:"required" omitempty:"true"`
}

type BudgetType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Budget struct {
	ID            uuid.UUID `json:"id"`
	Owner         string    `json:"owner"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Date          time.Time `json:"date"`
	Date_Start    time.Time `json:"date_start"`
	Date_End      time.Time `json:"date_end"`
	CurrencyCode  string    `json:"currency_code"`
	Base_Currency string    `json:"base_currency"`
	Is_Archived   bool      `json:"is_archived"`
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

type Currency struct {
	ID     int    `json:"id"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type BudgetInfo struct {
	BudgetName      string    `json:"budget_name"`
	OwnerName       string    `json:"owner_name"`
	OwnerEmail      string    `json:"owner_email"`
	BudgetType      string    `json:"budget_type"`
	CreatedAt       time.Time `json:"created_at"`
	BudgetStartDate time.Time `json:"budget_start_date"`
	BudgetEndDate   time.Time `json:"budget_end_date"`
	BudgetCurrency  string    `json:"budget_currency"`
	CurrencyCode    string    `json:"currency_code"`
	CurrencySymbol  string    `json:"currency_symbol"`
	IsArchived      bool      `json:"is_archived"`
}
