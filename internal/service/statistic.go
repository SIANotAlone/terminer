package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
)

type StatisticService struct {
	repo repository.Statistic
}

func NewStatisticService(repo repository.Statistic) *StatisticService {

	return &StatisticService{repo: repo}
}

func (s *StatisticService) GetProvidedDoneRecordsStatistic(user uuid.UUID, year int) (models.GaveDoneStatistic, error) {
	return s.repo.GetProvidedDoneRecordsStatistic(user, year)
}

func (s *StatisticService) GetProvidedRecordsProMonthStatistic(user uuid.UUID, year int) (models.RecordsProMonthStatistic, error) {
	var st models.RecordsProMonthStatistic
	Records, err := s.repo.GetProvidedRecordsProMonthStatistic(user, year)
	if err != nil {
		return st, err
	}
	st.Records = Records

	Services, err := s.repo.GetProvidedServicesProMonthStatistic(user, year)
	if err != nil {
		return st, err
	}
	st.Services = Services

	return st, nil
}

func (s *StatisticService) GetRecievedRecordsProMonthStatistic(user uuid.UUID, year int) (models.RecordsProMonthStatistic, error) {
	var st models.RecordsProMonthStatistic
	Records, err := s.repo.GetRecievedRecordsProMonthStatistic(user, year)
	if err != nil {
		return st, err
	}
	st.Records = Records

	Services, err := s.repo.GetRecievedServicesProMonthStatistic(user, year)
	if err != nil {
		return st, err
	}
	st.Services = Services

	return st, nil
}

func (s *StatisticService) GetMassageProTypeStatistic(user uuid.UUID, year int) (models.StatisticByType, error){
	st, err := s.repo.GetMassageProTypeStatistic(user, year)
	if err != nil {
		return models.StatisticByType{}, err
	}

	return fillByTypeStatistic(st), nil
}

func (s *StatisticService) GetResievedMassageProTypeStatistic(user uuid.UUID, year int) (models.StatisticByType, error){
	st, err := s.repo.GetResievedMassageProTypeStatistic(user, year)
	if err != nil {
		return models.StatisticByType{}, err
	}

	return fillByTypeStatistic(st), nil
}

func (s *StatisticService) GetResievedServicesTypes(user uuid.UUID, year int) (models.StatisticByType, error){
	st, err := s.repo.GetResievedServicesTypes(user, year)
	if err != nil {
		return models.StatisticByType{}, err
	}

	return fillByTypeStatistic(st), nil
}

func (s *StatisticService) GetProvidedServicesTypes(user uuid.UUID, year int) (models.StatisticByType, error){
	st, err := s.repo.GetProvidedServicesTypes(user, year)
	if err != nil {
		return models.StatisticByType{}, err
	}

	return fillByTypeStatistic(st), nil
}



func fillByTypeStatistic(st []models.ByType) models.StatisticByType {
	var bt models.StatisticByType

	for _, v := range st {
		bt.Types = append(bt.Types, v.Type)
		bt.Amount = append(bt.Amount, v.Total)
	}

	return bt
}

func (s *StatisticService) GetAvailableYears(user uuid.UUID) ([]int, error){
	return s.repo.GetAvailableYears(user)
}

func (s *StatisticService) GetMainStatistic(user uuid.UUID, year int) (models.MainStatistic, error){
	return s.repo.GetMainStatistic(user, year)
}