package repository

import (
	"database/sql"
	"fmt"
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	userTable  = "main.user"
	adminTable = "main.admin"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetAllUsers() ([]models.User, error) {

	query := fmt.Sprintf("SELECT uuid, first_name, last_name, country, email FROM %s", userTable)
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var users []models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Country, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil

}

func (r *UserPostgres) IsAdmin(id uuid.UUID) (bool, error) {
	query := fmt.Sprintf("SELECT user_id FROM %s WHERE user_id = $1", adminTable)
	row := r.db.QueryRow(query, id)
	var admin_id uuid.UUID
	if err := row.Scan(&admin_id); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
	}
	if admin_id == id {
		return true, nil
	}
	return false, nil
}
