package models

import (
	"time"

	"github.com/google/uuid"
)

type NewRecord struct {
	ServiceID       uuid.UUID `json:"service_id" binding:"required" omitempty:"true"`
	UserID          uuid.UUID `json:"-"`
	AvailableTimeID int       `json:"available_time_id" binding:"required" omitempty:"true"`
}

type DoneRecord struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}

type ServiceInfo struct {
	Name       string
	RecordDate time.Time
	TimeStart  time.Time
	TimeEnd    time.Time
}

type TerminsFromService struct {
	ID              uuid.UUID `json:"id"`
	Date            time.Time `json:"date"`
	Client          string    `json:"client"`
	Done            bool      `json:"done"`
	UserConfirm     bool      `json:"user_confirm"`
	AvailableTimeID int       `json:"available_time_id"`
	Time            time.Time `json:"time"`
	UserConfirmTime time.Time `json:"user_confirm_time"`
	DoneTime        time.Time `json:"done_time"`
	TimeStart       time.Time `json:"time_start"`
	TimeEnd         time.Time `json:"time_end"`
	Booked          bool      `json:"booked"`
	Service         string    `json:"service"`
	Description     string    `json:"description"`
}

type TerminsFromServiceResponce struct {
	Termins []TerminsFromService `json:"termins"`
	Total   int64                `json:"total"`
	Booked  int64                `json:"booked"`
}
type ServiceBookedInfo struct {
	Total  int64 `json:"total"`
	Booked int64 `json:"booked"`
}

type TerminsFromServiceInput struct {
	ServiceID uuid.UUID `json:"service_id" binding:"required" omitempty:"true"`
}
