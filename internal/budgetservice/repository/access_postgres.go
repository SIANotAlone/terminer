package repository

import (
	"terminer/internal/budgetservice/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccessPostgres struct {
	db *sqlx.DB
}

func NewAccessPostgres(db *sqlx.DB) *AccessPostgres {
	return &AccessPostgres{db: db}
}

func (r *AccessPostgres) ShareBudget(budgetID uuid.UUID, target_user uuid.UUID) error {
	query := `insert into budget.access(user_id, budget_id) 
values ($1, $2);`
	_, err := r.db.Exec(query, target_user, budgetID)
	return err
}

func (r *AccessPostgres) RevokeAccess(access_id uuid.UUID) error {
	query := `delete from budget.access where uuid = $1`
	_, err := r.db.Exec(query, access_id)
	return err
}

func (r *AccessPostgres) GetBudgetAccessList(budgetID uuid.UUID) ([]models.BudgetAccess, error) {
	query := `select dc.uuid, u.first_name || ' '|| u.last_name as user, u.email, dc.date from budget.access dc
left join main.user u on u.uuid = dc.user_id 
where dc.budget_id = $1;`
	row, err := r.db.Query(query, budgetID)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var accessList []models.BudgetAccess
	for row.Next() {
		var access models.BudgetAccess
		err := row.Scan(&access.ID, &access.UserName, &access.Email, &access.Date)
		if err != nil {
			return nil, err
		}
		accessList = append(accessList, access)
	}
	return accessList, nil
}

func (r *AccessPostgres) GetAllUsers() ([]models.User, error) {
	query := `select uuid, first_name || ' '|| last_name as name, email from main.user;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *AccessPostgres) GetBudgetOwnerID(budgetID uuid.UUID) (uuid.UUID, error) {
	query := `select owner_id from budget.budgets where uuid = $1;`
	var ownerID uuid.UUID
	err := r.db.Get(&ownerID, query, budgetID)
	return ownerID, err
}
