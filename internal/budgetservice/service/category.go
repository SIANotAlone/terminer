package service

import (
	"fmt"
	"terminer/internal/budgetservice/models"
	"terminer/internal/budgetservice/repository"

	"github.com/google/uuid"
)

type CategoryService struct {
	repo repository.Category
}

func NewCategoryService(repo repository.Category) *CategoryService {
	return &CategoryService{repo: repo}
}

func (s *CategoryService) CreateCategory(userID uuid.UUID, category models.NewCategory) (uuid.UUID, error) {
	return s.repo.CreateCategory(userID, category)
}

func (s *CategoryService) UpdateCategory(userID uuid.UUID, category models.UpdateCategory) error {
	CategoryOwner, err := s.repo.GetCategoryOwnerID(category.ID)
	if err != nil {
		return err
	}
	if CategoryOwner != userID {
		return fmt.Errorf("user is not the owner of the category")
	}
	return s.repo.UpdateCategory(category)
}

func (s *CategoryService) DeleteCategory(userID uuid.UUID, categoryID uuid.UUID) error {
	CategoryOwner, err := s.repo.GetCategoryOwnerID(categoryID)
	if err != nil {
		return err
	}
	if CategoryOwner != userID {
		return fmt.Errorf("user is not the owner of the category")
	}
	return s.repo.DeleteCategory(categoryID)
}

func (s *CategoryService) GetAvaliableCategories(userID uuid.UUID) ([]models.Category, error) {
	return s.repo.GetAvaliableCategories(userID)
}
