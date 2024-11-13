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
	CreateServiceType(models.ServiceType) error
	GetMyServices(user_id uuid.UUID) ([]models.MyService, error)
	GetAvailableService(user_id uuid.UUID) ([]models.AvailableService, error)
	GetAvailableTime(service_id uuid.UUID) ([]models.ServiceAvailableTime, error)

	GetUserTelegramID(user_id uuid.UUID) (string, error)
	GetAllUsersTelegramID() ([]string, error)
}

type Comment interface {
	CreateComment(comment models.Comment) (uuid.UUID, error)
	UpdateComment(comment models.UpdateComment) error
	DeleteComment(id uuid.UUID, user uuid.UUID) error
}

type Record interface {
	CreateRecord(record models.NewRecord) (uuid.UUID, error)
	DoneRecord(id uuid.UUID, user uuid.UUID) error
	ConfirmRecord(id uuid.UUID, user uuid.UUID) error
}

type Termin interface {
	GetAllPerformerTermins(user_id uuid.UUID) ([]models.Termin, error)
	GetAllUserTermins(user_id uuid.UUID) ([]models.Termin, error)
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
	Comment
	Termin
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Offering:      NewOfferingPostgres(db),
		User:          NewUserPostgres(db),
		Record:        NewRecordPostgres(db),
		Comment:       NewCommentPostgres(db),
		Termin:        NewTerminPostgres(db),
	}
}
