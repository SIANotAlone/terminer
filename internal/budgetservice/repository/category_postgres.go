package repository

import (
	"terminer/internal/budgetservice/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type CategoryPostgres struct {
	db *sqlx.DB
}

func NewCategoryPostgres(db *sqlx.DB) *CategoryPostgres {
	return &CategoryPostgres{db: db}
}

func (r *CategoryPostgres) CreateCategory(userID uuid.UUID, category models.NewCategory) (uuid.UUID, error) {
	query := `insert into budget.categories (name, description, user_id, type)
values ($1, $2, $3, $4) returning uuid;`
	var categoryID uuid.UUID
	err := r.db.Get(&categoryID, query, category.Name, category.Description, userID, category.Type)
	if err != nil {
		return uuid.Nil, err
	}
	return categoryID, nil
}

func (r *CategoryPostgres) UpdateCategory(category models.UpdateCategory) error {
	query := `update budget.categories set name = $1, description = $2, type = $3
where uuid = $4;`
	_, err := r.db.Exec(query, category.Name, category.Description, category.Type, category.ID)
	return err
}
func (r *CategoryPostgres) DeleteCategory(categoryID uuid.UUID) error {
	query := `delete from budget.categories where uuid = $1;`
	_, err := r.db.Exec(query, categoryID)
	return err
}

func (r *CategoryPostgres) GetAvaliableCategories(userID uuid.UUID) ([]models.Category, error) {
	query := `SELECT dc.uuid, dc.name, dc.description, u.first_name || ' '|| u.last_name as owner, dc.date, dc.is_based, dc.type
FROM budget.categories dc
LEFT JOIN main.user u on u.uuid = dc.user_id
WHERE dc.user_id = $1 OR dc.is_based = true
ORDER BY dc.type, dc.name;`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.Owner, &category.Created_At, &category.IsBased, &category.Type)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}
	return categories, nil
}

func (r *CategoryPostgres) GetCategoryOwnerID(categoryID uuid.UUID) (uuid.UUID, error) {
	var ownerID uuid.UUID
	query := `SELECT user_id FROM budget.categories WHERE uuid = $1;`
	err := r.db.Get(&ownerID, query, categoryID)
	if err != nil {
		return uuid.Nil, err
	}
	return ownerID, nil
}
