package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
)

type TerminService struct {
	repo repository.Termin
}

func NewTerminService(repo repository.Termin) *TerminService {
	
	return &TerminService{repo: repo}
}


func (s *TerminService) GetAllPerformerTermins(user_id uuid.UUID) ([]models.Termin, error){

	return s.repo.GetAllPerformerTermins(user_id)
}

func (s *TerminService) GetAllUserTermins(user_id uuid.UUID) ([]models.Termin, error){
	return s.repo.GetAllUserTermins(user_id)

}