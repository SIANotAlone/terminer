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