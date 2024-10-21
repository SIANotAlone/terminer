package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	UUID              uuid.UUID `json:"-"`
	Name              string    `json:"name" binding:"required" omitempty:"true"`
	Description       string    `json:"description" binding:"required" omitempty:"true"`
	Date              time.Time `json:"-"`
	DateEnd           time.Time `json:"date_end" binding:"required" omitempty:"true"`
	ServiceType       int       `json:"service_type" binding:"required" omitempty:"true"`
	PerformerID       uuid.UUID `json:"-"`
	Available_for_all bool      `json:"for_all" omitempty:"true"`
}

type Available_time struct {
	ID        int       `json:"-"`
	ServiceID uuid.UUID `json:"-`
	TimeStart string    `json:"time_start" binding:"required" omitempty:"true"`
	TimeEnd   string    `json:"time_end" binding:"required" omitempty:"true"`
}
type Available_for struct {
	ID        int       `json:"-"`
	ServiceID uuid.UUID `json:"-`
	UserID    uuid.UUID `json:"user_id" binding:"required" omitempty:"true"`
}

type NewService struct {
	Service        Service          `json:"service" binding:"required" omitempty:"true"`
	Available_time []Available_time `json:"available_time"`
	Available_for  []Available_for  `json:"available_for"`
}

type ServiceUpdate struct {
	UUID        *uuid.UUID `json:"id" binding:"required" omitempty:"true"`
	Name        *string    `json:"name" binding:"required" omitempty:"true"`
	Description *string    `json:"description" binding:"required" omitempty:"true"`
	DateEnd     *time.Time `json:"date_end" binding:"required" omitempty:"true"`
	ServiceType *int       `json:"service_type" binding:"required" omitempty:"true"`
}

type ServiceDelete struct {
	UUID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}

type ServiceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}