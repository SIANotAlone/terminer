package service

import (
	"fmt"
	"terminer/internal/models"
	"terminer/internal/observer"
	"terminer/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type RecordService struct {
	repo   repository.Record
	logger logrus.Logger
}

func NewRecordService(repo repository.Record, logger *logrus.Logger) *RecordService {
	return &RecordService{repo: repo, logger: *logger}
}

func (s *RecordService) CreateRecord(record models.NewRecord) (uuid.UUID, error) {

	repo, err := s.repo.CreateRecord(record)
	if err != nil {
		s.logger.Warn(err)
	}
	obs := observer.ConcreteObserver{}
	subject := observer.ConcreteSubject{}
	subject.Register(&obs)
	user_name, err := s.GetUserName(record.UserID)
	if err != nil {
		s.logger.Warn(err)
	}
	service_name, err := s.GetServiceName(record.ServiceID)
	if err != nil {
		s.logger.Warn(err)
	}
	message := fmt.Sprintf("Користувач __*%s*__ \nзаписався на вашу послугу: %s", user_name, service_name)
	owner_telegram_id, err := s.GetServiceOwnerTelegram(record.ServiceID)
	if err != nil {
		s.logger.Warn(err)
	}
	subject.Notify(owner_telegram_id, message)

	return repo, nil
}

func (s *RecordService) DoneRecord(id uuid.UUID, user uuid.UUID) error {
	return s.repo.DoneRecord(id, user)
}

func (s *RecordService) ConfirmRecord(id uuid.UUID, user uuid.UUID) error {
	return s.repo.ConfirmRecord(id, user)
}

func (s *RecordService) GetServiceOwnerTelegram(id uuid.UUID) (string, error) {
	return s.repo.GetServiceOwnerTelegram(id)
}

func (s *RecordService) GetUserName(user_id uuid.UUID) (string, error) {
	return s.repo.GetUserName(user_id)
}

func (s *RecordService) GetServiceName(id uuid.UUID) (string, error) {
	return s.repo.GetServiceName(id)
}
