package service

import (
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)

type TransactionService struct {
	repo repository.Transaction
}

func NewTransactionService(repo repository.Transaction) *TransactionService {
	return &TransactionService{repo: repo}
}

func (s *TransactionService) CreateTransaction(userID uuid.UUID, transaction models.NewTransaction) (uuid.UUID, error) {
	if transaction.GoalID != nil {
		return s.repo.CreateTransactionWithGoal(userID, transaction)
	}
	return s.repo.CreateTransactionWithoutGoal(userID, transaction)
}

func (s *TransactionService) UpdateTransaction(userID uuid.UUID, transaction models.UpdateTransaction) error {
	// 1. Получаем старую сумму для пересчета баланса (цели или бюджета)
	oldAmount, err := s.repo.GetTransactionAmountByID(transaction.TransactionID)
	if err != nil {
		return err
	}

	// 2. Просто обновляем. Репозиторий сам разберется с goal_id через RETURNING.
	return s.repo.UpdateTransaction(transaction, oldAmount)
}
func (s *TransactionService) DeleteTransaction(userID uuid.UUID, transactionID uuid.UUID) error {
	return s.repo.DeleteTransaction(transactionID)
}

func (s *TransactionService) GetTransactionsByBudget(userID uuid.UUID, budgetID uuid.UUID) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByBudget(budgetID)
}

func (s *TransactionService) GetBudgetIdByTransactionID(transactionID uuid.UUID) (uuid.UUID, error) {
	return s.repo.GetBudgetIdByTransactionID(transactionID)
}
