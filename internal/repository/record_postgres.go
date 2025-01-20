package repository

import (
	"fmt"
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

func (r *RecordPostgres) DoneRecord(id uuid.UUID, user uuid.UUID) error {
	query := fmt.Sprintf(`select u.uuid from main.record dc
			left join main.service s on s.uuid = dc.service_id
			left join main.user u on u.uuid = s.performer_id
			where dc.uuid = $1`)

	var owner uuid.UUID
	row := r.db.QueryRow(query, id)
	if err := row.Scan(&owner); err != nil {
		return err
	}
	if owner != user {
		return fmt.Errorf("user is not owner of service")
	}

	query2 := fmt.Sprintf(`UPDATE %s
			SET done = true, done_time = now()
			WHERE uuid = $1`, recordTable)
	_, err := r.db.Exec(query2, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *RecordPostgres) ConfirmRecord(id uuid.UUID, user uuid.UUID) error {
	query := fmt.Sprintf(`update %s set user_confirm = true, user_confirm_time = now() where uuid = $1 and user_id = $2`, recordTable)

	_, err := r.db.Exec(query, id, user)
	if err != nil {
		return err
	}

	return nil
}

func (r *RecordPostgres) GetServiceOwnerTelegram(service_id uuid.UUID) (string, error) {
	query := fmt.Sprintf(`SELECT u.telegram_chat_id  FROM main.service dc 
	LEFT JOIN main.user u on dc.performer_id = u.uuid
	WHERE dc.UUID = $1`)

	var chat_id string
	row := r.db.QueryRow(query, service_id)
	if err := row.Scan(&chat_id); err != nil {
		return "", err
	}
	return chat_id, nil
}

func (r *RecordPostgres) GetRecordOwnerTelegram(record_id uuid.UUID) (string, error) {
	query := fmt.Sprintf(`select u.telegram_chat_id from main.record dc
left join main.user u on u.uuid = dc.user_id
where dc.uuid = $1
`)

	var chat_id string
	row := r.db.QueryRow(query, record_id)
	if err := row.Scan(&chat_id); err != nil {
		return "", err
	}
	return chat_id, nil
}

func (r *RecordPostgres) GetUserName(user_id uuid.UUID) (string, error) {
	query := `select dc.first_name || ' ' || dc.last_name as name from main.user dc
	where dc.uuid  = $1`

	var name string
	row := r.db.QueryRow(query, user_id)
	if err := row.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

func (r *RecordPostgres) GetServiceName(service_id uuid.UUID) (string, error) {
	query := `select dc.name from main.service dc
	where dc.uuid = $1`

	var name string
	row := r.db.QueryRow(query, service_id)
	if err := row.Scan(&name); err != nil {
		return "", err
	}

	return name, nil
}

func (r *RecordPostgres) GetServiceInfo(service_id uuid.UUID) (models.ServiceInfo, error) {
	query := `select s.name, dc.date as rec_date, t.time_start, t.time_end from main.record dc
left join main.service s on dc.service_id = s.uuid
left join main.available_time t on t.id = dc.available_time_id
where dc.uuid = $1`

	var info models.ServiceInfo
	row := r.db.QueryRow(query, service_id)
	if err := row.Scan(&info.Name, &info.RecordDate, &info.TimeStart, &info.TimeEnd); err != nil {
		return models.ServiceInfo{}, err
	}

	return info, nil
}
