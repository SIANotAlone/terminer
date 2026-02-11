package models

import (
	"time"

	"github.com/google/uuid"
)

type NewCategory struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	UserID      string `*`
	Type        string `json:"type"`
}

type UpdateCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
}

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	Created_At  time.Time `json:"created_at"`
	IsBased     bool      `json:"is_based"`
	Type        string    `json:"type"`
}

type CategoryID struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}
