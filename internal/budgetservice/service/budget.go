package service

import (
	"fmt"
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)



type BudgetService struct {
	repo repository.Budget
}

func NewBudgetService(repo repository.Budget) *BudgetService {
	return &BudgetService{repo: repo}
}

func (s *BudgetService) CreateBudget(userID uuid.UUID, budget models.NewBudget) (uuid.UUID, error) {
	return s.repo.CreateBudget(userID, budget)
}

func (s *BudgetService) GetAvailableBudgets(userID uuid.UUID) ([]models.Budget, error) {
	return s.repo.GetAvailableBudgets(userID)
}
func (s *BudgetService) UpdateBudget(userID uuid.UUID, budget models.UpdateBudget) error {
	BudgetOwner , err := s.repo.GetBudgetOwnerID(budget.ID)
	if err != nil {
		return err
	}
	if BudgetOwner != userID {
		return fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.UpdateBudget(userID, budget)
}
func (s *BudgetService) DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	BudgetOwner , err := s.repo.GetBudgetOwnerID(budgetID)
	if err != nil {
		return err
	}
	if BudgetOwner != userID {
		return fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.DeleteBudget(userID, budgetID)
}
func (s *BudgetService) ArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	OwnerID , err := s.repo.GetBudgetOwnerID(budgetID)
	if err != nil {
		return err
	}
	if OwnerID != userID {
		return fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.ArchiveBudget(userID, budgetID)
}
func (s *BudgetService) UnArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error {
	OwnerID , err := s.repo.GetBudgetOwnerID(budgetID)
	if err != nil {
		return err
	}
	if OwnerID != userID {
		return fmt.Errorf("user is not the owner of the budget")
	}
	return s.repo.UnArchiveBudget(userID, budgetID)
}

func (s *BudgetService) GetBudgetTypes() ([]models.BudgetType, error) {
	return s.repo.GetBudgetTypes()
}

func (s *BudgetService) GetBudgetOwnerID(budgetID uuid.UUID) (uuid.UUID, error) {
	return s.repo.GetBudgetOwnerID(budgetID)
}