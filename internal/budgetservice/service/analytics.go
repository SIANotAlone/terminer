package service

import (
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)

type AnalyticsService struct {
	repo repository.Analytics
}

func NewAnalyticsService(repo repository.Analytics) *AnalyticsService {
	return &AnalyticsService{repo: repo}
}

func (s *AnalyticsService) GetDashboardData(budgetID, userID uuid.UUID) (*models.AnalyticsDashboard, error) {
	return s.repo.GetDashboardData(budgetID, userID)
}
