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

func (s *OfferingService) GetServiceOwner(id uuid.UUID) (uuid.UUID, error) {
	return s.repo.GetServiceOwner(id)
}

func (s *OfferingService) CreateServiceType(serviceType models.ServiceType) (error) {
	return s.repo.CreateServiceType(serviceType)
}

func (s *OfferingService) GetMyServices(user_id uuid.UUID) ([]models.MyService, error) {
	return s.repo.GetMyServices(user_id)
}

func (s *OfferingService) GetAvailableService(user_id uuid.UUID) ([]models.AvailableService, error) {
	return s.repo.GetAvailableService(user_id)
}