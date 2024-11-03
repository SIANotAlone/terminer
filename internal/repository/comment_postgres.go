package repository

import (
	"fmt"
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	commentTable = "main.comment"
)

type CommentPostgres struct {
	db *sqlx.DB
}

func NewCommentPostgres(db *sqlx.DB) *CommentPostgres {
	return &CommentPostgres{db: db}
}

func (r *CommentPostgres) CreateComment(comment models.Comment) (uuid.UUID, error) {
	query := fmt.Sprintf(`INSERT INTO %s (user_id, record_id, comment) VALUES ($1, $2, $3) RETURNING uuid`, commentTable)
	var id uuid.UUID
	err := r.db.QueryRow(query, comment.UserID, comment.RecordID, comment.Comment).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, nil
}

func (r *CommentPostgres) UpdateComment(comment models.UpdateComment) error {
	query := fmt.Sprintf(`UPDATE %s SET comment = $1, timechange=now() WHERE uuid = $2 and user_id = $3`, commentTable)
	_, err := r.db.Exec(query, comment.Comment, comment.ID, comment.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *CommentPostgres) DeleteComment(id uuid.UUID, user uuid.UUID) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE uuid = $1 and user_id = $2`, commentTable)
	row := r.db.QueryRow(query, id, user)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}
