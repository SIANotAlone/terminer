package models

import "github.com/google/uuid"


type NewRecord struct {
	ServiceID uuid.UUID `json:"service_id" binding:"required" omitempty:"true"` 
	UserID    uuid.UUID `json:"-"`
	AvailableTimeID int `json:"available_time_id" binding:"required" omitempty:"true"`
}

type DoneRecord struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}

