package service

import (
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)

type Authorization interface {
	ParseToken(accessToken string) (uuid.UUID, error)
}

type Budget interface {
	CreateBudget(userID uuid.UUID, budget models.NewBudget) (uuid.UUID, error)
	UpdateBudget(userID uuid.UUID, budget models.UpdateBudget) error
	DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error

	ArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error
	UnArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error

	GetAvailableBudgets(userID uuid.UUID, archived bool, limit int, offset int) ([]models.Budget, error)
	GetBudgetTypes() ([]models.BudgetType, error)
	GetBudgetOwnerID(budgetID uuid.UUID) (uuid.UUID, error)
	GetCurrencies() ([]models.Currency, error)
}

type Goal interface {
	CreateGoal(userID uuid.UUID, goal models.NewGoal) (uuid.UUID, error)
	UpdateGoal(userID uuid.UUID, goal models.UpdateGoal) error
	DeleteGoal(userID uuid.UUID, goalID uuid.UUID) error
	ArchiveGoal(userID uuid.UUID, goalID uuid.UUID) error
	UnArchiveGoal(userID uuid.UUID, goalID uuid.UUID) error

	GetAvailableGoals(userID uuid.UUID) ([]models.Goal, error)
	GetAllGoals(userID uuid.UUID) ([]models.Goal, error)
	GetGoalsTransactions(userID uuid.UUID, goalID uuid.UUID) ([]models.GoalTransaction, error)
}

type Transaction interface {
	CreateTransaction(userID uuid.UUID, transaction models.NewTransaction) (uuid.UUID, error)
	UpdateTransaction(userID uuid.UUID, transaction models.UpdateTransaction) error
	DeleteTransaction(userID uuid.UUID, transactionID uuid.UUID) error

	GetTransactionsByBudget(userID uuid.UUID, budgetID uuid.UUID) ([]models.Transaction, error)
	GetBudgetIdByTransactionID(transactionID uuid.UUID) (uuid.UUID, error)
}
type Category interface {
	CreateCategory(userID uuid.UUID, category models.NewCategory) (uuid.UUID, error)
	UpdateCategory(userID uuid.UUID, category models.UpdateCategory) error
	DeleteCategory(userID uuid.UUID, categoryID uuid.UUID) error

	GetAvaliableCategories(userID uuid.UUID) ([]models.Category, error)
}

type Access interface {
	ShareBudget(userID uuid.UUID, budgetID uuid.UUID, target_user uuid.UUID) (uuid.UUID, error)
	RevokeAccess(ownerID uuid.UUID, access_id uuid.UUID) error

	GetBudgetAccessList(userID uuid.UUID, budgetID uuid.UUID) ([]models.BudgetAccess, error)
	GetAllUsers() ([]models.User, error)

	HasUserAccessToBudget(userID uuid.UUID, budgetID uuid.UUID) (bool, error)
}

type Analytics interface {
	GetDashboardData(budgetID, userID uuid.UUID) (*models.AnalyticsDashboard, error)
}

type Service struct {
	Authorization Authorization
	Budget        Budget
	Goal          Goal
	Transaction   Transaction
	Category      Category
	Access        Access
	Analytics     Analytics
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(),
		Budget:        NewBudgetService(repos.Budget),
		Goal:          NewGoalService(repos.Goal),
		Transaction:   NewTransactionService(repos.Transaction),
		Category:      NewCategoryService(repos.Category),
		Access:        NewAccessService(repos.Access),
		Analytics:     NewAnalyticsService(repos.Analytics),
	}
}
