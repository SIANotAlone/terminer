package service

import (
	"fmt"
	"strconv"
	"terminer/internal/models"
	"terminer/internal/observer"
	"terminer/internal/repository"
	"time"

	"github.com/google/uuid"
)

type CommentService struct {
	repo repository.Comment
}

func NewCommentService(repo repository.Comment) *CommentService {
	return &CommentService{repo: repo}
}

func (s *CommentService) CreateComment(comment models.Comment) (uuid.UUID, error) {
	service_info, err := s.repo.GetServiceAndOwnerInfo(comment.RecordID)
	if err != nil {
		fmt.Println(err) // TODO: log
	}
	obs := observer.ConcreteObserver{}
	subject := observer.ConcreteSubject{}
	subject.Register(&obs)

	termin_date := s.getEscapedDate(service_info.TermineDate)
	if comment.UserID == service_info.PerformerID {
		message := fmt.Sprintf("Користувач __*%s*__ \nзалишив коментар під вашим записом на послугу: \"*%s*\" до нього\nПослуга від %s", service_info.PerformerName, service_info.ServiceName, termin_date)
		if service_info.RecordOwnerTG != "" {
			subject.Notify(service_info.RecordOwnerTG, message)
		}
	}
	if comment.UserID == service_info.RecordOwnerID {
		message := fmt.Sprintf("Користувач __*%s*__ \nзалишив коментар під записом на послугу \"*%s*\" до вас\nПослуга від %s", service_info.RecordOwnerName, service_info.ServiceName, termin_date)
		if service_info.PerformerTG != "" {
			subject.Notify(service_info.PerformerTG, message)
		}
	}

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

func (s *CommentService) getEscapedDate(date time.Time) string {

	str_date := "*" + strconv.Itoa(date.Day()) + "\\." + strconv.Itoa(int(date.Month())) + "\\." + strconv.Itoa(date.Year()) + "*"
	return str_date
}
