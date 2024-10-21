package repository

import (
	"fmt"
	"terminer/internal/models"

	"github.com/jmoiron/sqlx"
)

const (
	userTable = "main.user"
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
