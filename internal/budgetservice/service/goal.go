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
