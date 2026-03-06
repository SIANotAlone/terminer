package models


// AnalyticsDashboard — главная структура для ответа API
type AnalyticsDashboard struct {
	DonutChart   []CategoryExpense `json:"donut_chart"`
	BarChart     []PlanVsActual    `json:"bar_chart"`
	AreaChart    []DailyPulse      `json:"area_chart"`
	RadialChart  []GoalProgress    `json:"radial_chart"`
	TotalSpent   float64           `json:"total_spent"`    // Сумма за текущий период (€)
	PrevSpent    float64           `json:"prev_spent"`     // Сумма за прошлый период (€)
}

type CategoryExpense struct {
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
}

type PlanVsActual struct {
	CategoryName string  `json:"category_name"`
	Planned      float64 `json:"planned"`
	Actual       float64 `json:"actual"`
}

type DailyPulse struct {
	Date   string  `json:"date"` // Формат "02 Jan"
	Amount float64 `json:"amount"`
}

type GoalProgress struct {
	GoalName   string  `json:"goal_name"`
	Percentage float64 `json:"percentage"` // Рассчитывается как (current_saved / target_amount) * 100
}