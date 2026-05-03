package service

import (
	"fmt"
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)

type GoalService struct {
	repo repository.Goal
}

func NewGoalService(repo repository.Goal) *GoalService {
	return &GoalService{repo: repo}
}

func (s *GoalService) CreateGoal(userID uuid.UUID, goal models.NewGoal) (uuid.UUID, error) {
	return s.repo.CreateGoal(userID, goal)
}

func (s *GoalService) UpdateGoal(userID uuid.UUID, goal models.UpdateGoal) error {
	GoalOwner, err := s.repo.GetGoalOwnerID(goal.ID)
	if err != nil {
		return err
	}
	if GoalOwner != userID {
		return fmt.Errorf("user is not the owner of the goal")
	}
	return s.repo.UpdateGoal(goal)
}

func (s *GoalService) DeleteGoal(userID uuid.UUID, goalID uuid.UUID) error {
	GoalOwner, err := s.repo.GetGoalOwnerID(goalID)
	if err != nil {
		return err
	}
	if GoalOwner != userID {
		return fmt.Errorf("user is not the owner of the goal")
	}
	return s.repo.DeleteGoal(goalID)
}

func (s *GoalService) GetAvailableGoals(userID uuid.UUID) ([]models.Goal, error) {
	return s.repo.GetAvailableGoals(userID)
}

func (s *GoalService) GetAllGoals(userID uuid.UUID) ([]models.Goal, error) {
	return s.repo.GetAllGoals(userID)
}

func (s *GoalService) GetGoalsTransactions(userID uuid.UUID, goalID uuid.UUID) ([]models.GoalTransaction, error) {
	GoalOwner, err := s.repo.GetGoalOwnerID(goalID)
	if err != nil {
		return nil, err
	}
	if GoalOwner != userID {
		return nil, fmt.Errorf("user is not the owner of the goal")
	}
	return s.repo.GetGoalsTransactions(goalID)
}


func (s *GoalService) ArchiveGoal(userID uuid.UUID, goalID uuid.UUID) error  {
	return s.repo.ArchiveGoal(userID, goalID)
	
}

func (s *GoalService) UnArchiveGoal(userID uuid.UUID, goalID uuid.UUID) error  {
	return s.repo.UnArchiveGoal(userID, goalID)
}