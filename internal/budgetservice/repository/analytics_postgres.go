package repository

import (
	"database/sql"
	"fmt"
	"terminer/internal/budgetservice/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AnalyticsPostgres struct {
	db *sqlx.DB
}

func NewAnalyticsPostgres(db *sqlx.DB) *AnalyticsPostgres {
	return &AnalyticsPostgres{db: db}
}

func (r *AnalyticsPostgres) GetDashboardData(budgetID, userID uuid.UUID) (*models.AnalyticsDashboard, error) {
	datequery := `SELECT date_start, date_end FROM budget.budgets WHERE uuid = $1`
	var startDate, endDate time.Time
	err := r.db.QueryRow(datequery, budgetID).Scan(&startDate, &endDate)
	if err != nil {
		return nil, err
	}

	dashboard := &models.AnalyticsDashboard{
		DonutChart:  make([]models.CategoryExpense, 0),
		BarChart:    make([]models.PlanVsActual, 0),
		AreaChart:   make([]models.DailyPulse, 0),
		RadialChart: make([]models.GoalProgress, 0),
	}

	// 1. Donut Chart: Структура витрат (Фактические расходы по категориям)
	queryDonut := `
		SELECT c.name, SUM(t.amount) 
		FROM budget.transactions t
		JOIN budget.categories c ON t.category_id = c.uuid
		WHERE t.budget_id = $1 AND t.direction = 'EXPENSE' AND t.intent = 'ACTUAL'
		  
		GROUP BY c.name ORDER BY SUM(t.amount) DESC`

	rowsDonut, err := r.db.Query(queryDonut, budgetID)
	if err != nil {
		return nil, err
	}
	defer rowsDonut.Close()
	for rowsDonut.Next() {
		var ce models.CategoryExpense
		if err := rowsDonut.Scan(&ce.CategoryName, &ce.Amount); err != nil {
			return nil, err
		}
		dashboard.DonutChart = append(dashboard.DonutChart, ce)
	}

	// 2. Bar Chart: План vs Факт
	queryBar := `
		SELECT c.name,
		       COALESCE(SUM(CASE WHEN t.intent = 'PLANNED' THEN t.amount END), 0) as planned,
		       COALESCE(SUM(CASE WHEN t.intent = 'ACTUAL' THEN t.amount END), 0) as actual
		FROM budget.transactions t
		JOIN budget.categories c ON t.category_id = c.uuid
		WHERE t.budget_id = $1 AND t.direction = 'EXPENSE'
		GROUP BY c.name`

	rowsBar, err := r.db.Query(queryBar, budgetID)
	if err != nil {
		return nil, err
	}
	defer rowsBar.Close()
	for rowsBar.Next() {
		var pva models.PlanVsActual
		if err := rowsBar.Scan(&pva.CategoryName, &pva.Planned, &pva.Actual); err != nil {
			return nil, err
		}
		dashboard.BarChart = append(dashboard.BarChart, pva)
	}

	// 3. Area Chart: Пульс витрат (Группировка по дням)
	queryArea := `
		SELECT DATE(t.date) as day, SUM(t.amount) 
		FROM budget.transactions t
		WHERE t.budget_id = $1 AND t.direction = 'EXPENSE' AND t.intent = 'ACTUAL'
		
		GROUP BY DATE(t.date) ORDER BY day ASC`

	rowsArea, err := r.db.Query(queryArea, budgetID)
	if err != nil {
		return nil, err
	}
	var ukrMonthsShort = []string{
		"Січня", "Лютого", "Березня", "Квітня", "Травня", "Червня",
		"Липня", "Серпня", "Вересня", "Жовтня", "Листопада", "Грудня",
	}
	defer rowsArea.Close()
	for rowsArea.Next() {
		var dp models.DailyPulse
		var date time.Time
		if err := rowsArea.Scan(&date, &dp.Amount); err != nil {
			return nil, err
		}

		// Отримуємо номер місяця (від 1 до 12) і беремо назву з масиву
		monthIdx := int(date.Month()) - 1
		dp.Date = fmt.Sprintf("%02d %s", date.Day(), ukrMonthsShort[monthIdx])

		dashboard.AreaChart = append(dashboard.AreaChart, dp)
	}

	// 4. Radial Chart: Прогресс целей (привязаны к юзеру)
	queryGoals := `
		SELECT target_name, target_amount, current_saved 
		FROM budget.accumulation_goals 
		WHERE user_id = $1`

	rowsGoals, err := r.db.Query(queryGoals, userID)
	if err != nil {
		return nil, err
	}
	defer rowsGoals.Close()
	for rowsGoals.Next() {
		var g models.GoalProgress
		var targetAmount, currentSaved float64
		if err := rowsGoals.Scan(&g.GoalName, &targetAmount, &currentSaved); err != nil {
			return nil, err
		}

		if targetAmount > 0 {
			g.Percentage = (currentSaved / targetAmount) * 100
		}
		dashboard.RadialChart = append(dashboard.RadialChart, g)
	}

	// 5. Статистика (Текущий месяц vs Прошлый месяц) для пустого стейта с советом
	duration := endDate.Sub(startDate)
	prevEndDate := startDate.Add(-time.Second)
	prevStartDate := prevEndDate.Add(-duration)

	// ШАГ 1: Находим ID предыдущего бюджета по вычисленным датам.
	// ВАЖНО: Замените 'owner_id' на то поле, которое объединяет бюджеты одного пользователя/категории
	// (например, user_id, account_id или group_id), чтобы не захватить чужой бюджет на эти же даты.
	var prevBudgetID *string // Используем указатель, на случай если прошлого бюджета не существует

	queryPrevBudget := `
        SELECT uuid 
        FROM budget.budgets 
        WHERE owner_id = (SELECT owner_id FROM budget.budgets WHERE uuid = $1)
          AND date_start = $2 
          AND date_end = $3
        LIMIT 1`

	err = r.db.QueryRow(queryPrevBudget, budgetID, prevStartDate, prevEndDate).Scan(&prevBudgetID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// ШАГ 2: Считаем суммы, привязываясь ТОЛЬКО к ID бюджетов.
	// Даты самих транзакций здесь больше не играют роли.
	queryStats := `
        SELECT 
            COALESCE(SUM(CASE WHEN budget_id = $1 THEN amount ELSE 0 END), 0) as current_spent,
            COALESCE(SUM(CASE WHEN budget_id = $2 THEN amount ELSE 0 END), 0) as prev_spent
        FROM budget.transactions
        WHERE budget_id IN ($1, $2) 
          AND direction = 'EXPENSE' 
          AND intent = 'ACTUAL'`

	err = r.db.QueryRow(queryStats, budgetID, prevBudgetID).
		Scan(&dashboard.TotalSpent, &dashboard.PrevSpent)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return dashboard, nil
}
