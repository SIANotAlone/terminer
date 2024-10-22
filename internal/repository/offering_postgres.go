package repository

import (
	"fmt"
	"terminer/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	servicesTable       = "main.service"
	availableForTable   = "main.available_for"
	available_timeTable = "main.available_time"
	service_typeTable   = "main.service_type"
	time_layout         = "15:04"
)

type OfferingPostgres struct {
	db *sqlx.DB
}

func NewOfferingPostgres(db *sqlx.DB) *OfferingPostgres {
	return &OfferingPostgres{db: db}
}

func (r *OfferingPostgres) CreateOffering(offering models.NewService) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	var id uuid.UUID
	create_service_query := fmt.Sprintf("INSERT INTO %s (name, description, date, date_end, service_type_id, performer_id, available_for_all) VALUES ($1, $2, CURRENT_DATE, $3, $4, $5, $6) RETURNING uuid", servicesTable)
	row := tx.QueryRow(create_service_query, offering.Service.Name, offering.Service.Description, offering.Service.DateEnd, offering.Service.ServiceType, offering.Service.PerformerID, offering.Service.Available_for_all)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return id, err
	}
	for _, value := range offering.Available_for {
		create_avalable_for_query := fmt.Sprintf("INSERT INTO %s (service_id, user_id) VALUES ($1, $2)", availableForTable)
		_, err := tx.Exec(create_avalable_for_query, id, value.UserID)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
	}
	for _, value := range offering.Available_time {
		create_available_time_query := fmt.Sprintf("INSERT INTO %s (service_id, time_start, time_end) VALUES ($1, $2, $3)", available_timeTable)
		time_start, err := time.Parse(time_layout, value.TimeStart)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		time_end, err := time.Parse(time_layout, value.TimeEnd)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		_, err = tx.Exec(create_available_time_query, id, time_start.Format(time_layout), time_end.Format(time_layout))
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
	}

	return id, tx.Commit()
}

func (r *OfferingPostgres) UpdateService(service models.ServiceUpdate) error {

	query := fmt.Sprintf("UPDATE %s SET name = $1, description = $2, date_end = $3, service_type_id = $4 WHERE uuid = $5", servicesTable)
	row := r.db.QueryRow(query, service.Name, service.Description, service.DateEnd, service.ServiceType, service.UUID)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) DeleteService(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uuid = $1", servicesTable)
	row := r.db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) GetTypes() ([]models.ServiceType, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s", service_typeTable)
	row, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var service_types []models.ServiceType
	for row.Next() {
		var service_type models.ServiceType
		if err := row.Scan(&service_type.ID, &service_type.Name); err != nil {
			return nil, err
		}
		service_types = append(service_types, service_type)
	}
	return service_types, nil
}

func (r *OfferingPostgres) GetServiceOwner(id uuid.UUID) (uuid.UUID, error) {
	query := fmt.Sprintf("SELECT performer_id FROM %s WHERE uuid = $1", servicesTable)
	row := r.db.QueryRow(query, id)
	var owner_id uuid.UUID
	if err := row.Scan(&owner_id); err != nil {
		return owner_id, err
	}
	return owner_id, nil
}

func (r *OfferingPostgres) CreateServiceType(s models.ServiceType) error {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1)", service_typeTable)
	_, err := r.db.Exec(query, s.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) GetMyServices(user_id uuid.UUID) ([]models.MyService, error) {
	query := `select dc.uuid, dc.name, dc.description, dc.date, dc.date_end, st.name as service_type from main.service dc
	left join main.user u on dc.performer_id = u.uuid
	left join main.service_type st on dc.service_type_id = st.id 
	where dc.performer_id = $1 and date_end < CURRENT_DATE`
	var myservices []models.MyService
	row, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var ms models.MyService
		if err := row.Scan(&ms.ID, &ms.Name, &ms.Description, &ms.Date, &ms.DateEnd, &ms.ServiceType); err != nil {
			return nil, err
		}
		myservices = append(myservices, ms)
	}
	

	return myservices, nil
}
