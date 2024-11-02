package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
)
type RecordService struct {
	repo repository.Record
}

func NewRecordService(repo repository.Record) *RecordService {
	return &RecordService{repo: repo}
}

func (s *RecordService) CreateRecord(record models.NewRecord) (uuid.UUID, error) {
	return s.repo.CreateRecord(record)
}

func (s *RecordService) DoneRecord(id uuid.UUID, user uuid.UUID) (error) {
	return s.repo.DoneRecord(id, user)
}

func (s *RecordService) ConfirmRecord(id uuid.UUID, user uuid.UUID) (error) {
	return s.repo.ConfirmRecord(id, user)
}