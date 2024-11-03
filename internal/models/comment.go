package models

import "github.com/google/uuid"

type Comment struct {
	UserID    uuid.UUID `-`
	RecordID  uuid.UUID `json:"record_id"`
	Comment   string    `json:"comment"`
}

type UpdateComment struct {
	ID      uuid.UUID `json:"id" binding:"required" omitempty:"true"`
	UserID  uuid.UUID `-`
	Comment string `json:"comment"`
}

type DeleteComment struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}