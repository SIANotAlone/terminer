package service

import (
	"fmt"
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)

type AccessService struct {
	repo repository.Access
}

func NewAccessService(repo repository.Access) *AccessService {
	return &AccessService{repo: repo}
}

func (s *AccessService) ShareBudget(userID uuid.UUID, budgetID uuid.UUID, target_user uuid.UUID) error {
	owner, err := s.repo.GetBudgetOwnerID(budgetID)
	if err != nil {
		return err
	}
	if owner != userID {
		return fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.ShareBudget(budgetID, target_user)
}
func (s *AccessService) RevokeAccess(ownerID uuid.UUID, access_id uuid.UUID) error {
	owner, err := s.repo.GetBudgetOwnerID(access_id)
	if err != nil {
		return err
	}
	if owner != ownerID {
		return fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.RevokeAccess(access_id)
}

func (s *AccessService) GetBudgetAccessList(userID uuid.UUID, budgetID uuid.UUID) ([]models.BudgetAccess, error) {
	owner, err := s.repo.GetBudgetOwnerID(budgetID)
	if err != nil {
		return nil, err
	}
	if owner != userID {
		return nil, fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.GetBudgetAccessList(budgetID)
}

func (s *AccessService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}
