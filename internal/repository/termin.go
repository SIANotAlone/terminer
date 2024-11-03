package repository

import (
	"fmt"
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const serviceTable = "main.service"

type TerminPostgres struct {
	db *sqlx.DB
}

func NewTerminPostgres(db *sqlx.DB) *TerminPostgres {
	return &TerminPostgres{db: db}
}

func (r *TerminPostgres) GetAllPerformerTermins(user_id uuid.UUID) ([]models.Termin, error) {
	var performerTermins []models.Termin
	query := fmt.Sprintf(`select r.uuid as record_id,
		u.first_name || ' ' || u.last_name as performer,
		st.name as type,  dc.name as service,
		dc.description, dc.date, dc.date_end, 
		r.time as record_time, a_t.time_start, a_t.time_end
		from %s dc
		left join main.record r on r.service_id = dc.uuid
		left join main.service_type st on st.id = dc.service_type_id
		left join main.user u on u.uuid = dc.performer_id
		left join main.available_time a_t on a_t.id = r.available_time_id
		where performer_id = $1 
		and r.done != true 
		and r.user_confirm != true `, serviceTable)

	row, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var performerTermin models.Termin
		if err := row.Scan(&performerTermin.RecordID, &performerTermin.Performer, &performerTermin.Type, &performerTermin.Service,
			&performerTermin.Description, &performerTermin.Date, &performerTermin.DateEnd, &performerTermin.RecordTime, &performerTermin.TimeStart, &performerTermin.TimeEnd); err != nil {
			return nil, err
		}
		performerTermins = append(performerTermins, performerTermin)
	}

	return performerTermins, nil
}

func (r *TerminPostgres) GetAllUserTermins(user_id uuid.UUID) ([]models.Termin, error) {
	var userTermins []models.Termin
	query := fmt.Sprintf(`select dc.uuid as record_id, u.first_name || ' ' || u.last_name as performer,
		st.name as type, s.name as service, s.description, s.date, s.date_end,
		dc.time as record_time, 
		a_t.time_start, a_t.time_end
		from %s dc
		left join main.service s on s.uuid = dc.service_id 
		left join main.user u on u.uuid = s.performer_id
		left join main.service_type st on st.id = s.service_type_id
		left join main.available_time a_t on a_t.id = dc.available_time_id
		where user_id = $1 and dc.done=false`, recordTable)

	row, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}

	for row.Next() {
		var userTermin models.Termin
		if err := row.Scan(&userTermin.RecordID, &userTermin.Performer, &userTermin.Type, &userTermin.Service,
			&userTermin.Description, &userTermin.Date, &userTermin.DateEnd, &userTermin.RecordTime, &userTermin.TimeStart, &userTermin.TimeEnd); err != nil {
			return nil, err
		}
		userTermins = append(userTermins, userTermin)
	}
	return userTermins, nil
}
