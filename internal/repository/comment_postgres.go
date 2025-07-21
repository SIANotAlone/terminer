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
	query := fmt.Sprintf(`UPDATE %s SET comment = $1, timestampchange=now() WHERE uuid = $2 and user_id = $3`, commentTable)
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

func (r *CommentPostgres) GetCommentsOnRecord(record_id uuid.UUID, user_id uuid.UUID) (models.CommentsList, error) {
	var comments models.CommentsList
	query := `select dc.uuid, u.first_name||' '||u.last_name as comment_owner, dc.comment, dc.timestamp, dc.timestampchange,
		-- **************************
		-- чи є користувач власником коментаря
		CASE 
			WHEN dc.user_id = $1 THEN true
			ELSE false
			END as is_my_comment
		-- **************************

		from main.comment dc
		left join main.user u on u.uuid = dc.user_id 

		where dc.record_id=$2
		`

	row, err := r.db.Query(query, user_id, record_id)
	if err != nil {
		return comments, err
	}
	for row.Next() {
		var comment models.CommentOnRecord
		if err := row.Scan(&comment.ID, &comment.CommentOwner, &comment.Comment, &comment.Created, &comment.Updated, &comment.IsMyComment); err != nil {
			return comments, err
		}
		comments.CommentsList = append(comments.CommentsList, comment)
	}

	return comments, nil
}

func (r *CommentPostgres) GetTerminsWithComments(record_id uuid.UUID) ([]models.TerminsWithComments, error) {
	var twc []models.TerminsWithComments // termins with comments
	query := `SELECT r.uuid, max(s.name) as service_name, r.date, r.done, r.user_confirm, r.time, r.user_confirm_time, r.done_time
	FROM main.comment dc
	LEFT JOIN main.record r ON r.uuid = dc.record_id
	LEFT JOIN main.service s ON s.uuid = r.service_id
	WHERE 
		s.performer_id = $1
		OR r.user_id = $1
	GROUP BY r.uuid
	`

	row, err := r.db.Query(query, record_id)
	if err != nil {
		return twc, err
	}
	for row.Next() {
		var list models.TerminsWithComments
		if err := row.Scan(&list.ID, &list.ServiceName, &list.Created, &list.Done, &list.UserConfirm, &list.Created_time, &list.User_confirm_time, &list.Done_time); err != nil {
			return twc, err
		}

		twc = append(twc, list)
	
	}
	

	return twc, nil
}