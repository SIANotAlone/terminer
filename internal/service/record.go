package service

import (
	"fmt"
	"strings"
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

	var isBooked bool
	isBooked, err := s.repo.CheckAvailableTime(record.AvailableTimeID, record.ServiceID)
	if err != nil {
		s.logger.Warn(err)
	}
	if isBooked == true {
		return uuid.Nil, fmt.Errorf("time is not available")
	}
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
	message := fmt.Sprintf("Користувач __*%s*__ \nзаписався на вашу послугу: %s", escapeMarkdownV2(user_name), escapeMarkdownV2(service_name))
	owner_telegram_id, err := s.GetServiceOwnerTelegram(record.ServiceID)
	if err != nil {
		s.logger.Warn(err)
	}
	subject.Notify(owner_telegram_id, message)

	return repo, nil
}

func (s *RecordService) DoneRecord(id uuid.UUID, user uuid.UUID) error {
	err := s.repo.DoneRecord(id, user)
	if err != nil {
		return err
	}
	obs := observer.ConcreteObserver{}
	subject := observer.ConcreteSubject{}
	subject.Register(&obs)
	user_name, err := s.GetUserName(user)
	if err != nil {
		s.logger.Warn(err)
	}
	service_info, err := s.GetServiceInfo(id)
	if err != nil {
		s.logger.Warn(err)
	}
	owner_telegram_id, err := s.GetRecordOwnerTelegram(id)
	if err != nil {
		return err
	}
	record_date := service_info.RecordDate.Format("02.01.2006")
	time_start := service_info.TimeStart.Format("15:04")
	time_end := service_info.TimeEnd.Format("15:04")
	message := fmt.Sprintf("__*%s*__ позначив ваш запис на послугу: __*%s*__  \nЗапис від %s\nЗапис в проміжку між %s та %s\nПозначено як: __*Виконано*__\n",
		escapeMarkdownV2(user_name), escapeMarkdownV2(service_info.Name),
		escapeMarkdownV2(record_date), escapeMarkdownV2(time_start), escapeMarkdownV2(time_end))
	subject.Notify(owner_telegram_id, message)
	return nil
}

func (s *RecordService) ConfirmRecord(id uuid.UUID, user uuid.UUID) error {
	return s.repo.ConfirmRecord(id, user)
}

func (s *RecordService) GetServiceOwnerTelegram(id uuid.UUID) (string, error) {
	return s.repo.GetServiceOwnerTelegram(id)
}
func (s *RecordService) GetRecordOwnerTelegram(record_id uuid.UUID) (string, error) {
	return s.repo.GetRecordOwnerTelegram(record_id)
}

func (s *RecordService) GetUserName(user_id uuid.UUID) (string, error) {
	return s.repo.GetUserName(user_id)
}

func (s *RecordService) GetServiceName(id uuid.UUID) (string, error) {
	return s.repo.GetServiceName(id)
}

func (s *RecordService) GetServiceInfo(record_id uuid.UUID) (models.ServiceInfo, error) {
	return s.repo.GetServiceInfo(record_id)
}

func (s *RecordService) GetTerminsFromService(service_id uuid.UUID) (models.TerminsFromServiceResponce, error) {
	var responce models.TerminsFromServiceResponce
	var termins []models.TerminsFromService
	termins, err := s.repo.GetTerminsFromService(service_id)
	if err != nil {
		return models.TerminsFromServiceResponce{}, err
	}
	responce.Termins = termins
	var service_booked_info models.ServiceBookedInfo
	service_booked_info, err = s.repo.GetServiceBookedInfo(service_id)
	if err != nil {
		return models.TerminsFromServiceResponce{}, err
	}
	responce.Booked = service_booked_info.Booked
	responce.Total = service_booked_info.Total
	return responce, nil

}
func escapeMarkdownV2(input string) string {
	replacer := strings.NewReplacer(
		".", "\\.",
		"-", "\\-",
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		"!", "\\!",
	)
	return replacer.Replace(input)
}
