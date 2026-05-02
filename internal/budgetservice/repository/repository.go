package repository

import (
	"terminer/internal/budgetservice/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Budget interface {
	CreateBudget(userID uuid.UUID, budget models.NewBudget) (uuid.UUID, error)
	UpdateBudget(userID uuid.UUID, budget models.UpdateBudget) error
	DeleteBudget(userID uuid.UUID, budgetID uuid.UUID) error
	ArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error
	UnArchiveBudget(userID uuid.UUID, budgetID uuid.UUID) error

	GetAvailableBudgets(userID uuid.UUID) ([]models.Budget, error)
	GetAvailableBudgetsWithArchived(userID uuid.UUID, limit int, offset int) ([]models.Budget, error)

	GetBudgetTypes() ([]models.BudgetType, error)

	GetBudgetOwnerID(budgetID uuid.UUID) (uuid.UUID, error)
	GetCurrencies() ([]models.Currency, error)
}

type Goal interface {
	CreateGoal(userID uuid.UUID, goal models.NewGoal) (uuid.UUID, error)
	UpdateGoal(goal models.UpdateGoal) error
	DeleteGoal(goalID uuid.UUID) error
	ArchiveGoal(userID uuid.UUID, goalID uuid.UUID) error
	UnArchiveGoal(userID uuid.UUID, goalID uuid.UUID) error

	GetAvailableGoals(userID uuid.UUID) ([]models.Goal, error)
	GetAllGoals(userID uuid.UUID) ([]models.Goal, error)
	GetGoalOwnerID(goalID uuid.UUID) (uuid.UUID, error)
	GetGoalsTransactions(goalID uuid.UUID) ([]models.GoalTransaction, error)
}

type Transaction interface {
	CreateTransactionWithGoal(userID uuid.UUID, transaction models.NewTransaction) (uuid.UUID, error)
	CreateTransactionWithoutGoal(userID uuid.UUID, transaction models.NewTransaction) (uuid.UUID, error)
	UpdateTransaction(userID uuid.UUID, newTx models.UpdateTransaction) error
	DeleteTransaction(transactionID uuid.UUID) error

	GetTransactionsByBudget(budgetID uuid.UUID) ([]models.Transaction, error)
	GetTransactionAmountByID(transactionID uuid.UUID) (decimal.Decimal, error)

	HasGoal(TransactionID uuid.UUID) (bool, error)
	GetTrasactionWihtGoalAmount(transactionID uuid.UUID) (decimal.Decimal, uuid.UUID, error)

	GetBudgetIdByTransactionID(transactionID uuid.UUID) (uuid.UUID, error)
}
type Category interface {
	CreateCategory(userID uuid.UUID, category models.NewCategory) (uuid.UUID, error)
	UpdateCategory(category models.UpdateCategory) error
	DeleteCategory(categoryID uuid.UUID) error

	GetAvaliableCategories(userID uuid.UUID) ([]models.Category, error)
	GetCategoryOwnerID(categoryID uuid.UUID) (uuid.UUID, error)
}

type Access interface {
	ShareBudget(budgetID uuid.UUID, target_user uuid.UUID) (uuid.UUID, error)
	RevokeAccess(access_id uuid.UUID) error

	GetBudgetAccessList(budgetID uuid.UUID) ([]models.BudgetAccess, error)
	GetAllUsers() ([]models.User, error)

	GetBudgetOwnerID(budgetID uuid.UUID) (uuid.UUID, error)

	HasUserAccessToBudget(userID uuid.UUID, budgetID uuid.UUID) (bool, error)
	GetBudgetIDByAccessID(accessID uuid.UUID) (uuid.UUID, error)
}

type Analytics interface {
	GetDashboardData(budgetID, userID uuid.UUID) (*models.AnalyticsDashboard, error)
}

type Repository struct {
	Budget      Budget
	Goal        Goal
	Transaction Transaction
	Category    Category
	Access      Access
	Analytics   Analytics
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Budget:      NewBudgetPostgres(db),
		Goal:        NewGoalPostgres(db),
		Transaction: NewTransactionPostgres(db),
		Category:    NewCategoryPostgres(db),
		Access:      NewAccessPostgres(db),
		Analytics:   NewAnalyticsPostgres(db),
	}
}
