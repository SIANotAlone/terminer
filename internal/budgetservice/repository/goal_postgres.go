package repository

import (
	"terminer/internal/budgetservice/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type GoalPostgres struct {
	db *sqlx.DB
}

func NewGoalPostgres(db *sqlx.DB) *GoalPostgres {
	return &GoalPostgres{db: db}
}

func (r *GoalPostgres) CreateGoal(userID uuid.UUID, goal models.NewGoal) (uuid.UUID, error) {
	query := `insert into budget.accumulation_goals (user_id,target_name, target_amount, target_date, currency_id)
values ($1, $2, $3, $4, $5);`
	_, err := r.db.Exec(query, userID, goal.TargetName, goal.TargetAmount, goal.TargetDate, goal.CurrencyID)
	if err != nil {
		return uuid.Nil, err
	}
	return uuid.New(), nil
}

func (r *GoalPostgres) UpdateGoal(goal models.UpdateGoal) error {
	query := `update budget.accumulation_goals SET target_name=$1, target_amount=$2, target_date=$3, currency_id = $4
where uuid = $5;`
	_, err := r.db.Exec(query, goal.TargetName, goal.TargetAmount, goal.TargetDate, goal.CurrencyID, goal.ID)
	if err != nil {
		return err
	}
	return nil
}

func (r *GoalPostgres) DeleteGoal(goalID uuid.UUID) error {
	query := `delete from budget.accumulation_goals where uuid = $1;`
	_, err := r.db.Exec(query, goalID)
	if err != nil {
		return err
	}
	return nil
}

func (r *GoalPostgres) GetAvailableGoals(userID uuid.UUID) ([]models.Goal, error) {
	query := `SELECT dc.uuid, dc.target_name, dc.target_amount, dc.target_date, dc.current_saved, c.code, c.name
FROM budget.accumulation_goals dc
LEFT JOIN budget.currencies c on c.id = dc.currency_id 
WHERE dc.user_id = $1;`

	var goals []models.Goal
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var goal models.Goal
		err := rows.Scan(&goal.ID, &goal.TargetName, &goal.TargetAmount, &goal.TargetDate, &goal.CurrentSaved, &goal.CurrencyCode, &goal.CurrencyName)
		if err != nil {
			return nil, err
		}
		monthsNeeded := monthsFromNowSimple(goal.TargetDate)
		if monthsNeeded <= 0 {
			goal.RequiredMonthlySave = goal.TargetAmount
		} else {
			goal.RequiredMonthlySave = goal.TargetAmount.Sub(goal.CurrentSaved).Div(decimal.NewFromInt(int64(monthsNeeded)))
		}
		goals = append(goals, goal)
	}
	return goals, nil
}

func monthsFromNowSimple(to time.Time) int {
	now := time.Now()
	y1, m1, _ := now.Date()
	y2, m2, _ := to.Date()

	return (y2-y1)*12 + int(m2-m1)
}

func (r *GoalPostgres) GetGoalOwnerID(goalID uuid.UUID) (uuid.UUID, error) {
	var ownerID uuid.UUID
	query := `SELECT user_id FROM budget.accumulation_goals WHERE uuid = $1;`
	err := r.db.Get(&ownerID, query, goalID)
	if err != nil {
		return uuid.Nil, err
	}
	return ownerID, nil
}
