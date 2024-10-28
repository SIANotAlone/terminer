package repository

import (
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user models.UserRegistration) (uuid.UUID, error)
	GetUser(email string, password string) (uuid.UUID, error)
}

type Offering interface {
	CreateOffering(offering models.NewService) (uuid.UUID, error)
	UpdateService(models.ServiceUpdate) error
	DeleteService(id uuid.UUID) error
	GetTypes() ([]models.ServiceType, error)
	GetServiceOwner(id uuid.UUID) (uuid.UUID, error)
	CreateServiceType(models.ServiceType) (error)
	GetMyServices(user_id uuid.UUID) ([]models.MyService, error)
	GetAvailableService(user_id uuid.UUID) ([]models.AvailableService, error)
	GetAvailableTime(service_id uuid.UUID) ([]models.ServiceAvailableTime, error)
}
type Record interface {
	CreateRecord(record models.NewRecord) (uuid.UUID, error)
}


type User interface {
	GetAllUsers() ([]models.User, error)
	IsAdmin(id uuid.UUID) (bool, error)
}

type Repository struct {
	Authorization
	Offering
	User
	Record
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Offering:      NewOfferingPostgres(db),
		User:          NewUserPostgres(db),
		Record:        NewRecordPostgres(db),
	}
}
