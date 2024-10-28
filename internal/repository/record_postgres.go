package repository

import (
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	recordTable = "main.record"
)

type RecordPostgres struct {
	db *sqlx.DB
}

func NewRecordPostgres(db *sqlx.DB) *RecordPostgres {
	return &RecordPostgres{db: db}
}

func (r *RecordPostgres) CreateRecord(record models.NewRecord) (uuid.UUID, error) {
	query := `insert into main.record (service_id, user_id, available_time_id)
		values ($1, $2, $3) returning uuid`

	var record_id uuid.UUID
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	row := tx.QueryRow(query, record.ServiceID, record.UserID, record.AvailableTimeID)
	if err := row.Scan(&record_id); err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}
	query2 := `UPDATE main.available_time 
			SET booked = true 
			WHERE id = $1`

	_, err = tx.Exec(query2, record.AvailableTimeID)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	tx.Commit()
	return record_id, nil
}
