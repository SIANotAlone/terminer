package models

import (
	"time"

	"github.com/google/uuid"
)

type Termin struct {
	RecordID     uuid.UUID `json:"record_id"`
	Performer    string    `json:"performer"`
	Type         string    `json:"type"`
	Service      string    `json:"service"`
	Description  string    `json:"description"`
	Date         time.Time `json:"date"`
	DateEnd      time.Time `json:"date_end"`
	RecordTime   time.Time `json:"record_time"`
	TimeStart    string    `json:"time_start"`
	TimeEnd      string    `json:"time_end"`
	Done         bool      `json:"done"`
	User_confirm bool      `json:"user_confirm"`
}
