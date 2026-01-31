package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	UserID   uuid.UUID `-`
	RecordID uuid.UUID `json:"record_id"`
	Comment  string    `json:"comment"`
}

type UpdateComment struct {
	ID      uuid.UUID `json:"id" binding:"required" omitempty:"true"`
	UserID  uuid.UUID `-`
	Comment string    `json:"comment"`
}

type DeleteComment struct {
	ID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}

type CommentsList struct {
	CommentsList []CommentOnRecord `json:"comments_list"`
}

type CommentOnRecord struct {
	ID           uuid.UUID `json:"id"`
	CommentOwner string    `json:"comment_owner"`
	Comment      string    `json:"comment"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	IsMyComment  bool      `json:"is_my_comment"`
}

type GetCommentsInput struct {
	RecordID uuid.UUID `json:"record_id" binding:"required" omitempty:"true"`
}

type TerminsWithComments struct {
	ID                uuid.UUID `json:"id"`
	ServiceName       string    `json:"service"`
	Created           time.Time `json:"created"`
	Done              bool      `json:"done"`
	UserConfirm       bool      `json:"user_confirm"`
	Created_time      time.Time `json:"created_time"`
	User_confirm_time time.Time `json:"user_confirm_time"`
	Done_time         time.Time `json:"done_time"`
}

type ServiceAndOwnerInfo struct {
	ServiceID       uuid.UUID
	ServiceName     string
	ServiceDate     time.Time
	TermineDate     time.Time
	RecordOwnerID   uuid.UUID
	RecordOwnerName string
	RecordOwnerTG   string
	PerformerID     uuid.UUID
	PerformerName   string
	PerformerTG     string
}
