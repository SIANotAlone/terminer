package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
)

type OfferingService struct {
	repo repository.Offering
}

func NewOfferingService(repo repository.Offering) *OfferingService {
	return &OfferingService{repo: repo}
}

func (s *OfferingService) CreateService(offering models.NewService) (uuid.UUID, error) {
	return s.repo.CreateOffering(offering)
}

func (s *OfferingService) UpdateService(service models.ServiceUpdate) error {
	return s.repo.UpdateService(service)
}

func (s *OfferingService) DeleteService(id uuid.UUID) error {
	return s.repo.DeleteService(id)
}

func (s *OfferingService) GetTypes() ([]models.ServiceType, error) {
	return s.repo.GetTypes()
}