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
	if transaction.CategoryID != uuid.Nil {
		old_amount, err := s.repo.GetTransactionAmountByID(transaction.TransactionID)
		if err != nil {
			return err
		}
		return s.repo.UpdateTransactionWithGoal(transaction, old_amount)
	}
	return s.repo.UpdateTransactionWithoutGoal(transaction)
}

func (s *TransactionService) DeleteTransaction(userID uuid.UUID, transactionID uuid.UUID) error {
	hasGoal, err := s.repo.TransactionHasGoal(transactionID)
	if err != nil {
		return err
	}
	if hasGoal {
		return s.repo.DeleteTransactionWithGoal(transactionID)
	}
	return s.repo.DeleteTransactionWithoutGoal(transactionID)
}

func (s *TransactionService) GetTransactionsByBudget(userID uuid.UUID, budgetID uuid.UUID) ([]models.Transaction, error) {
	return s.repo.GetTransactionsByBudget(budgetID)
}
