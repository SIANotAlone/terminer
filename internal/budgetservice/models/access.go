package models

import (
	"time"

	"github.com/google/uuid"
)

type BudgetAccess struct {
	ID       uuid.UUID `json:"id"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	Date     time.Time `json:"date"`
}

type ShareBudgetInput struct {
	BudgetID   uuid.UUID `json:"budget_id" binding:"required"`
	TargetUser uuid.UUID `json:"target_user" binding:"required"`
}

type RevokeAccessInput struct {
	AccessID uuid.UUID `json:"access_id" binding:"required"`
}
