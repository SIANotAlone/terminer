package repository

import (
	"fmt"
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	usersTable = "main.user"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user models.UserRegistration) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (email, hash, first_name, last_name, date_birth, country, registration_date) VALUES ($1, $2, $3, $4, $5, $6, current_date) RETURNING uuid", usersTable)
	row := r.db.QueryRow(query, user.Email, user.Password, user.FirstName, user.LastName, user.DateOfBirth, user.Country)
	if err := row.Scan(&id); err != nil {
		return id, err
	}

	return id, nil
}

func(r *AuthPostgres) GetUser(email string, password string) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("SELECT uuid FROM %s WHERE email = $1 AND hash = $2", usersTable)
	err := r.db.Get(&id, query, email, password)
	return id, err
}