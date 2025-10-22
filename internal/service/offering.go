package service

import (
	"fmt"
	"strconv"
	"terminer/internal/models"
	"terminer/internal/observer"
	"terminer/internal/repository"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type OfferingService struct {
	repo   repository.Offering
	logger logrus.Logger
}

func NewOfferingService(repo repository.Offering, logger *logrus.Logger) *OfferingService {
	return &OfferingService{repo: repo, logger: *logger}
}

func (s *OfferingService) CreateService(offering models.NewService) (uuid.UUID, error) {
	repo, err := s.repo.CreateOffering(offering)
	if err != nil {
		s.logger.Warn(err)
	}
	obs := observer.ConcreteObserver{}
	subject := observer.ConcreteSubject{}
	subject.Register(&obs)
	var available_time string
	for _, time := range offering.Available_time {
		available_time += fmt.Sprintf("\n З *%s* по *%s*", time.TimeStart, time.TimeEnd)
	}

	message := fmt.Sprintf("Для *Вас* доступна нова __*послуга*__😍 \nНазва: %s\nОпис: %s\nПослуга доступна до: %s\nПослуга доступна в проміжку: %s",
		escapeMarkdownV2(offering.Service.Name), escapeMarkdownV2(offering.Service.Description),
		s.getEscapedDate(offering.Service.DateEnd), available_time)

	if offering.Service.Available_for_all == true {
		s.notificateAllUsers(&subject, message, offering.Service.PerformerID)
		return repo, nil
	} else {
		s.notificate_available_for_users(&subject, message, offering.Available_for)
	}
	return repo, nil
}

func (s *OfferingService) CreatePromoService(offering models.NewPromoService) (models.PromocodeServiceInfo, error) {
	repo, err := s.repo.CreatePromoService(offering)
	if err != nil {
		s.logger.Warn(err)
	}
	id := repo
	info, err := s.repo.GetServicePromocode(id)
	if err != nil {
		s.logger.Warn(err)
	}
	return info, nil
}

func (s OfferingService) ValidatePromoCode(code string) (models.PromocodeValidation, error) {
	var validation models.PromocodeValidation
	info, err := s.GetPromoCodeInfo(code)
	if err != nil {
		return validation, err
	}
	if info.Available_for > 0 {
		validation.Valid = false
	} else {
		validation.Valid = true
		validation.PromeService = info

	}
	return validation, nil
}

func (s OfferingService) ActivatePromoCode(code string, user_id uuid.UUID) error {
	info, err := s.GetPromoCodeInfo(code)
	if err != nil {
		return err
	}
	author_telegram_id, err := s.repo.GetPromoCodeServiceOwnerTelegramID(code)
	if err != nil {
		s.logger.Warn(err)
	}
	if author_telegram_id == "" {
		s.logger.Warn("service owner telegram id is empty")
		return s.repo.ActivatePromoCode(info.Service_ID, user_id)
	}
	service_info, err := s.repo.GetPromoCodeInfo(code)
	if err != nil {
		s.logger.Warn(err)
	}

	user_name, err := s.repo.GetUserName(user_id)
	if err != nil {
		s.logger.Warn(err)
	}

	message := fmt.Sprintf("Ваш промокод на послугу: \"*%s*\" \nбуло активовано користувачем __*%s*__", service_info.Name, user_name)
	obs := observer.ConcreteObserver{}
	subject := observer.ConcreteSubject{}
	subject.Register(&obs)
	s.notificate_one_user(&subject, message, author_telegram_id)
	// TODO:
	// У майбутньому скоріше за все необхідно буде додати очищення промокоду, бо вони можут співпадати
	return s.repo.ActivatePromoCode(info.Service_ID, user_id)
}

func (s *OfferingService) GetMyActualServices(user_id uuid.UUID) ([]models.MyActualService, error) {
	return s.repo.GetMyActualServices(user_id)
}

func (s *OfferingService) GetHistoryMyServices(user_id uuid.UUID, limit int64, offset int64) (models.UserServiceHistory, error) {
	var history models.UserServiceHistory
	h, err := s.repo.GetHistoryMyServices(user_id, limit, offset)
	if err != nil {
		s.logger.Warn(err)
	}
	history.History = h
	t, err := s.repo.GetTotalUserServices(user_id)
	if err != nil {
		s.logger.Warn(err)
	}
	history.Total = t
	return history, nil
}

func (s *OfferingService) GetPromoCodeInfo(code string) (models.PromocodeInfo, error) {
	return s.repo.GetPromoCodeInfo(code)
}

func (s *OfferingService) UpdateService(service models.ServiceUpdate) error {
	return s.repo.UpdateService(service)
}

func (s *OfferingService) DeleteService(id uuid.UUID, user uuid.UUID) error {
	owner, err := s.repo.GetServiceOwner(id)
	if err != nil {
		s.logger.Warn(err)
	}
	if owner != user {

		return fmt.Errorf("you are not owner of this service")
	}

	info, err := s.repo.GetFullServiceInfo(id)
	if err != nil {
		s.logger.Warn(err)
	}
	for _, value := range info.Available_time_Info {
		if value.Booked == true {
			return fmt.Errorf("this service has already booked appointments(termins)")
		}
	}

	return s.repo.DeleteService(id)

}

func (s *OfferingService) GetTypes() ([]models.ServiceType, error) {
	return s.repo.GetTypes()
}

func (s *OfferingService) GetMassageTypes() ([]models.MassageType, error) {
	return s.repo.GetMassageTypes()
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

func (s *OfferingService) GetFullServiceInfo(id uuid.UUID) (models.FullServiceInformation, error) {
	return s.repo.GetFullServiceInfo(id)
}

func (s *OfferingService) notificateAllUsers(subject *observer.ConcreteSubject, message string, exeption uuid.UUID) {
	users, err := s.repo.GetAllUsersTelegramID()
	if err != nil {
		s.logger.Warn(err)
	}
	// виключаємо автора послуги зі списку сповіщень
	exeption_telegram_id, err := s.repo.GetUserTelegramID(exeption)
	if err != nil {
		s.logger.Warn(err)
	}
	for _, user := range users {
		if user != exeption_telegram_id {
			subject.Notify(user, message)
		}
	}
}
func (s *OfferingService) notificate_one_user(subject *observer.ConcreteSubject, message string, tg_id string) {
	subject.Notify(tg_id, message)
}

func (s *OfferingService) EditService(service models.ServiceInformation, user uuid.UUID) error {
	owner, err := s.repo.GetServiceOwner(service.ID)
	if err != nil {
		s.logger.Warn(err)
	}
	if owner == user {
		return s.repo.EditService(service)
	} else {
		return fmt.Errorf("you are not owner of this service")
	}

}

func (s *OfferingService) NewAvailableTime(at models.NewAvailableTime, user uuid.UUID) (int, error) {
	owner, err := s.repo.GetServiceOwner(at.ServiceID)
	if err != nil {
		s.logger.Warn(err)
	}
	if owner == user {
		return s.repo.NewAvailableTime(at)
	} else {
		return 0, fmt.Errorf("you are not owner of this service")
	}
}

func (s *OfferingService) DeleteAvailableTime(at models.DeleteAvailableTime, user uuid.UUID) error {
	return s.repo.DeleteAvailableTime(at, user)
}

func (s *OfferingService) NewAvailableFor(af models.NewAvailableFor, user uuid.UUID) (int, error) {
	owner, err := s.repo.GetServiceOwner(af.ServiceID)
	if err != nil {
		s.logger.Warn(err)
	}
	if owner == user {
		return s.repo.NewAvailableFor(af)
	} else {
		return 0, fmt.Errorf("you are not owner of this service")
	}
}

func (s *OfferingService) DeleteAvailableFor(af models.DeleteAvailableFor, user uuid.UUID) error {
	return s.repo.DeleteAvailableFor(af, user)
}

func (s *OfferingService) notificate_available_for_users(subject *observer.ConcreteSubject, message string, users []models.Available_for) {
	for _, user := range users {
		tg_id, err := s.repo.GetUserTelegramID(user.UserID)
		if err != nil {
			s.logger.Warn(err)
		}
		subject.Notify(tg_id, message)

	}

}

func (s *OfferingService) getEscapedDate(date time.Time) string {

	str_date := "*" + strconv.Itoa(date.Day()) + "\\." + strconv.Itoa(int(date.Month())) + "\\." + strconv.Itoa(date.Year()) + "*"
	return str_date
}
