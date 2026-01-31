package service

import (
	"terminer/internal/mainservice/models"
	"terminer/internal/mainservice/repository"

	"github.com/google/uuid"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) IsAdmin(id uuid.UUID) (bool, error) {
	return s.repo.IsAdmin(id)
}