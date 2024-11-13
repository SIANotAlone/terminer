package service

import (
	"fmt"
	"terminer/internal/models"
	"terminer/internal/observer"
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
	repo, err := s.repo.CreateOffering(offering)
	if err != nil {
		fmt.Println(err)
	}
	obs := observer.ConcreteObserver{}
	subject := observer.ConcreteSubject{}
	subject.Register(&obs)
	message := fmt.Sprintf("Для вас доступна нова послуга... \n%s\n%s",
		offering.Service.Name, offering.Service.Description)

	if offering.Service.Available_for_all == true {
		s.notificateAllUsers(&subject, message)
		return repo, nil
	} else {
		s.notificate_available_for_users(&subject, message, offering.Available_for)
	}

	return repo, nil
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

func (s *OfferingService) CreateServiceType(serviceType models.ServiceType) error {
	return s.repo.CreateServiceType(serviceType)
}

func (s *OfferingService) GetMyServices(user_id uuid.UUID) ([]models.MyService, error) {
	return s.repo.GetMyServices(user_id)
}

func (s *OfferingService) GetAvailableService(user_id uuid.UUID) ([]models.AvailableService, error) {
	return s.repo.GetAvailableService(user_id)
}

func (s *OfferingService) GetAvailableTime(service_id uuid.UUID) ([]models.ServiceAvailableTime, error) {
	return s.repo.GetAvailableTime(service_id)
}

func (s *OfferingService) GetUserTelegramID(user_id uuid.UUID) (string, error) {
	return s.repo.GetUserTelegramID(user_id)
}

func (s *OfferingService) GetAllUsersTelegramID() ([]string, error) {
	return s.repo.GetAllUsersTelegramID()
}

func (s *OfferingService) notificateAllUsers(subject *observer.ConcreteSubject, message string) {
	users, err := s.repo.GetAllUsersTelegramID()
	if err != nil {
		println(err)
	}
	for _, user := range users {
		subject.Notify(user, message)
	}
}

func (s *OfferingService) notificate_available_for_users(subject *observer.ConcreteSubject, message string, users []models.Available_for) {
	for _, user := range users {
		tg_id, err := s.repo.GetUserTelegramID(user.UserID)
		if err != nil {
			println(err)
		}
		subject.Notify(tg_id, message)

	}

}
