package service

import (
	"terminer/internal/models"
	"terminer/internal/repository"

	"github.com/google/uuid"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment models.Comment) (uuid.UUID, error) {
	return s.repo.CreateComment(comment)
}


func (s *CommentService) UpdateComment(comment models.UpdateComment) error {
	return s.repo.UpdateComment(comment)
}

func (s *CommentService) DeleteComment(id uuid.UUID, user uuid.UUID) error {
	return s.repo.DeleteComment(id, user)
}