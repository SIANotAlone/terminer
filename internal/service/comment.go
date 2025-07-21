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

func (s *CommentService) GetCommentsOnRecord(record_id uuid.UUID, user_id uuid.UUID) (models.CommentsList, error) {
	return s.repo.GetCommentsOnRecord(record_id, user_id)
}

func (s *CommentService) GetTerminsWithComments(record_id uuid.UUID) ([]models.TerminsWithComments, error) {
	return s.repo.GetTerminsWithComments(record_id)
}