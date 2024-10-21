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
}

type User interface {
	GetAllUsers() ([]models.User, error)
}

type Repository struct {
	Authorization
	Offering
	User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Offering:      NewOfferingPostgres(db),
		User:          NewUserPostgres(db),
	}
}
