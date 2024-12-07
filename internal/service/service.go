package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	GetTypes() ([]models.ServiceType, error)
	GetServiceOwner(id uuid.UUID) (uuid.UUID, error)
	CreateServiceType(models.ServiceType) error
	GetMyServices(user_id uuid.UUID) ([]models.MyService, error)
	GetAvailableService(user_id uuid.UUID) ([]models.AvailableService, error)
	GetAvailableTime(service_id uuid.UUID) ([]models.ServiceAvailableTime, error)

	GetUserTelegramID(user_id uuid.UUID) (string, error)
	GetAllUsersTelegramID() ([]string, error)
}

type Record interface {
	CreateRecord(record models.NewRecord) (uuid.UUID, error)
	DoneRecord(id uuid.UUID, user uuid.UUID) error
	ConfirmRecord(id uuid.UUID, user uuid.UUID) error

	GetServiceOwnerTelegram(id uuid.UUID) (string, error)
	GetUserName(user_id uuid.UUID) (string, error)
	GetServiceName(id uuid.UUID) (string, error)
}

type Comment interface {
	CreateComment(comment models.Comment) (uuid.UUID, error)
	UpdateComment(comment models.UpdateComment) error
	DeleteComment(id uuid.UUID, user uuid.UUID) error
}

type Termin interface {
	GetAllPerformerTermins(user_id uuid.UUID) ([]models.Termin, error)
	GetAllUserTermins(user_id uuid.UUID) ([]models.Termin, error)
}

type User interface {
	GetAllUsers() ([]models.User, error)
	IsAdmin(id uuid.UUID) (bool, error)
}

type Service struct {
	Authorization
	Offering
	User
	Record
	Comment
	Termin
	logger logrus.Logger
}

func NewService(repos *repository.Repository, logger *logrus.Logger) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		Offering:      NewOfferingService(repos.Offering, logger),
		User:          NewUserService(repos.User),
		Record:        NewRecordService(repos.Record, logger),
		Comment:       NewCommentService(repos.Comment),
		Termin:        NewTerminService(repos.Termin),
		logger:        *logger,
	}
}
