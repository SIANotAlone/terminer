package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	CreateUser(user models.UserRegistration) (uuid.UUID, error)
	GenerateToken(email string, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}

type Offering interface {
	CreateService(offering models.NewService) (uuid.UUID, error)
	UpdateService(models.ServiceUpdate) error
	DeleteService(id uuid.UUID) error
}
type User interface {
	GetAllUsers() ([]models.User, error)
}

type Service struct {
	Authorization
	Offering
	User
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Offering:      NewOfferingService(repos.Offering),
		User:          NewUserService(repos.User),
	}
}
