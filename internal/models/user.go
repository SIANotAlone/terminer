package models

import "github.com/google/uuid"

type UserRegistration struct {
	FirstName      string `json:"first_name" binding:"required" omitempty:"true"`
	LastName       string `json:"last_name" binding:"required" omitempty:"true"`
	DateOfBirth    string `json:"date_of_birth" binding:"required" omitempty:"true"`
	Country        string `json:"country"`
	Email          string `json:"email" binding:"required" omitempty:"true"`
	Password       string `json:"password" binding:"required" omitempty:"true"`
	TelegramChatID string `json:"telegram_chat_id"`
}

type UserSignIn struct {
	Email    string `json:"email" binding:"required" omitempty:"true"`
	Password string `json:"password" binding:"required" omitempty:"true"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Country   string    `json:"country"`
	Email     string    `json:"email"`
}
