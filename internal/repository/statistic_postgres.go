package repository

import (
	"terminer/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type StatisticPostgres struct {
	db *sqlx.DB
}

func NewStatisticPostgres(db *sqlx.DB) *StatisticPostgres {
	return &StatisticPostgres{db: db}
}

func (r *StatisticPostgres) GetProvidedDoneRecordsStatistic(user uuid.UUID, year int) (models.GaveDoneStatistic, error) {
	var st models.GaveDoneStatistic
	query := `WITH vars AS (
    SELECT 
        $1::uuid AS user_id,
        $2::int  AS year
)
SELECT 
    (
        -- кол-во предоставленных услуг
        SELECT COUNT(*) 
        FROM main.record dc
        LEFT JOIN main.service s ON dc.service_id = s.uuid
        LEFT JOIN main.user u ON u.uuid = s.performer_id
        WHERE s.performer_id = vars.user_id
          AND EXTRACT(YEAR FROM dc.date) = vars.year AND dc.done=true
    ) AS gave,
    (
        -- кол-во полученных услуг
        SELECT COUNT(*) 
        FROM main.record dc
        WHERE dc.user_id = vars.user_id
          AND EXTRACT(YEAR FROM dc.date) = vars.year AND dc.done=true
    ) AS got,
    (
        -- имя пользователя
        SELECT first_name || ' ' || last_name
        FROM main.user 
        WHERE uuid = vars.user_id
    ) AS user_name
FROM vars;
`

	row := r.db.QueryRow(query, user, year)
	if err := row.Scan(&st.Gave, &st.Got, &st.UserName); err != nil {
		return st, err
	}
	return st, nil
}

func (r *StatisticPostgres) GetProvidedRecordsProMonthStatistic(user uuid.UUID, year int) (models.ProMonth, error) {
	var st models.ProMonth

	query := `
WITH vars AS (
  SELECT $1::int AS year, $2::uuid AS user_id
),
months AS (
  SELECT generate_series(
           make_date((SELECT year FROM vars), 1, 1),
           make_date((SELECT year FROM vars), 12, 1),
           interval '1 month'
         ) AS month_start
),
counts AS (
  SELECT date_trunc('month', dc.date) AS month_start, COUNT(*) AS gave 
  FROM main.record dc
  LEFT JOIN main.service s ON dc.service_id = s.uuid
  WHERE s.performer_id = (SELECT user_id FROM vars)
    AND EXTRACT(YEAR FROM dc.date) = (SELECT year FROM vars)
  GROUP BY month_start
)
SELECT
  COALESCE(c.gave, 0) AS gave
FROM months m
LEFT JOIN counts c USING (month_start)
ORDER BY m.month_start;
`

	var gave int
	var monthes []int

	rows, err := r.db.Query(query, year, user)
	if err != nil {
		return st, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&gave); err != nil {
			return st, err
		}
		monthes = append(monthes, gave)
	}

	st = fillStatistics(monthes)
	return st, nil
}

func (r *StatisticPostgres) GetProvidedServicesProMonthStatistic(user uuid.UUID, year int) (models.ProMonth, error) {
	var st models.ProMonth

	query := `-- Услуг создано
WITH params AS (
  SELECT $1::int AS year -- здесь меняете год
),
months AS (
  SELECT generate_series(
    make_date((SELECT year FROM params), 1, 1),
    make_date((SELECT year FROM params), 12, 1),
    interval '1 month'
  ) AS month_start
)
SELECT
 -- m.month_start,
  -- to_char(m.month_start, 'FMMonth') AS month_name,
  COALESCE((
    SELECT COUNT(*)
    FROM main.service dc
    WHERE dc.performer_id = $2
      AND dc.date >= m.month_start
      AND dc.date <  m.month_start + interval '1 month'
  ), 0) AS gave
FROM months m
ORDER BY m.month_start;`
	var gave int
	var monthes []int

	rows, err := r.db.Query(query, year, user)
	if err != nil {
		return st, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&gave); err != nil {
			return st, err
		}
		monthes = append(monthes, gave)
	}

	st = fillStatistics(monthes)
	return st, nil

}

func (r *StatisticPostgres) GetRecievedRecordsProMonthStatistic(user uuid.UUID, year int) (models.ProMonth, error) {
	var st models.ProMonth
	query := `WITH vars AS (
  SELECT 
    $1::int AS year,
    $2::uuid AS user_id
),
months AS (
  SELECT generate_series(
           make_date((SELECT year FROM vars), 1, 1),
           make_date((SELECT year FROM vars), 12, 1),
           interval '1 month'
         ) AS month_start
),
counts AS (
  SELECT date_trunc('month', dc.date) AS month_start,
         COUNT(*) AS got
  FROM main.record dc
  WHERE dc.user_id = (SELECT user_id FROM vars)
    AND EXTRACT(YEAR FROM dc.date) = (SELECT year FROM vars)
  GROUP BY month_start
)
SELECT
  -- m.month_start,
  -- to_char(m.month_start, 'FMMonth') AS month_name,
  COALESCE(c.got, 0) AS got
FROM months m
LEFT JOIN counts c USING (month_start)
ORDER BY m.month_start;

`
	var gave int
	var monthes []int

	rows, err := r.db.Query(query, year, user)
	if err != nil {
		return st, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&gave); err != nil {
			return st, err
		}
		monthes = append(monthes, gave)
	}

	st = fillStatistics(monthes)
	return st, nil

}

func (r *StatisticPostgres) GetRecievedServicesProMonthStatistic(user uuid.UUID, year int) (models.ProMonth, error) {
	var st models.ProMonth
	query := `--Услуг получено
WITH params AS (
  SELECT $1::int AS year -- меняете здесь год
),
months AS (
  SELECT generate_series(
    make_date((SELECT year FROM params), 1, 1),
    make_date((SELECT year FROM params), 12, 1),
    interval '1 month'
  ) AS month_start
),
counts AS (
  SELECT
    date_trunc('month', dc.date) AS month_start,
    COUNT(DISTINCT s.uuid) AS services_count
  FROM main.record dc
  LEFT JOIN main.service s ON s.uuid = dc.service_id
  WHERE dc.user_id = $2
    AND EXTRACT(YEAR FROM dc.date) = (SELECT year FROM params)
  GROUP BY date_trunc('month', dc.date)
)
SELECT
  -- m.month_start,
  -- to_char(m.month_start, 'FMMonth') AS month_name,
  COALESCE(c.services_count, 0) AS services_count
FROM months m
LEFT JOIN counts c USING (month_start)
ORDER BY m.month_start;`
	var gave int
	var monthes []int

	rows, err := r.db.Query(query, year, user)
	if err != nil {
		return st, err
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&gave); err != nil {
			return st, err
		}
		monthes = append(monthes, gave)
	}

	st = fillStatistics(monthes)
	return st, nil

}

func fillStatistics(st []int) models.ProMonth {
	var pm models.ProMonth
	pm.January = st[0]
	pm.February = st[1]
	pm.March = st[2]
	pm.April = st[3]
	pm.May = st[4]
	pm.June = st[5]
	pm.July = st[6]
	pm.August = st[7]
	pm.September = st[8]
	pm.October = st[9]
	pm.November = st[10]
	pm.December = st[11]
	return pm
}

func (r *StatisticPostgres) GetMassageProTypeStatistic(user uuid.UUID, year int) ([]models.ByType, error) {
	var st []models.ByType
	query := `select 
    mt.name,
    count(r.*) as total
from main.massage_type mt
left join main.service s 
    on mt.id = s.massage_type_id
left join main.record r 
    on r.service_id = s.uuid
   and extract(year from r.date) = $1
   and r.done = true
   and s.performer_id = $2
group by mt.id, mt.name
order by mt.id;`

	rows, err := r.db.Query(query, year, user)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var bt models.ByType
		if err := rows.Scan(&bt.Type, &bt.Total); err != nil {
			return st, err
		}
		st = append(st, bt)
	}
	return st, nil

}

func (r *StatisticPostgres) GetResievedMassageProTypeStatistic(user uuid.UUID, year int) ([]models.ByType, error) {
	var st []models.ByType
	query := `select 
    mt.name,
    count(r.*) as total
from main.massage_type mt
left join main.service s 
    on mt.id = s.massage_type_id
left join main.record r 
    on r.service_id = s.uuid
   and extract(year from r.date) = $1
   and r.done = true
   and r.user_id = $2
group by mt.id, mt.name
order by mt.id;`
	rows, err := r.db.Query(query, year, user)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var bt models.ByType
		if err := rows.Scan(&bt.Type, &bt.Total); err != nil {
			return st, err
		}
		st = append(st, bt)
	}
	return st, nil
}

func (r *StatisticPostgres) GetResievedServicesTypes(user uuid.UUID, year int) ([]models.ByType, error) {
	var st []models.ByType
	query := `select max(dc.name), count(*) from main.service_type dc

left join main.service s on dc.id = s.service_type_id
left join main.record r on r.service_id = s.uuid

where r.user_id = $1
    AND EXTRACT(YEAR FROM s.date) = $2 and r.done =true

group by dc.id`

	rows, err := r.db.Query(query, user, year)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var bt models.ByType
		if err := rows.Scan(&bt.Type, &bt.Total); err != nil {
			return st, err
		}
		st = append(st, bt)
	}
	return st, nil
}

func (r *StatisticPostgres) GetProvidedServicesTypes(user uuid.UUID, year int) ([]models.ByType, error) {
	var st []models.ByType
	query := `select max(dc.name), count(*) from main.service_type dc

left join main.service s on dc.id = s.service_type_id
left join main.record r on r.service_id = s.uuid

where s.performer_id = $1
    AND EXTRACT(YEAR FROM s.date) = $2 and r.done =true

group by dc.id`

	rows, err := r.db.Query(query, user, year)
	if err != nil {
		return st, err
	}
	defer rows.Close()
	for rows.Next() {
		var bt models.ByType
		if err := rows.Scan(&bt.Type, &bt.Total); err != nil {
			return st, err
		}
		st = append(st, bt)
	}
	return st, nil
}

func (r *StatisticPostgres) GetAvailableYears(user uuid.UUID) ([]int, error) {
	var years []int
	query := `SELECT DISTINCT EXTRACT(YEAR FROM date)::int AS year
FROM main.service
where performer_id = $1
ORDER BY year;
;
`

	rows, err := r.db.Query(query, user)
	if err != nil {
		return years, err
	}
	defer rows.Close()
	for rows.Next() {
		var year int
		if err := rows.Scan(&year); err != nil {
			return years, err
		}
		years = append(years, year)
	}
	return years, nil
}

func (r *StatisticPostgres) GetMainStatistic(user uuid.UUID, year int) (models.MainStatistic, error) {
	var st models.MainStatistic
	query := `
WITH vars AS (
    SELECT $1::uuid AS uid,$2::int AS year
)
SELECT 
	-- Услуг созданно пользователем
    (SELECT count(*) FROM main.service s WHERE s.performer_id = vars.uid  
	AND EXTRACT(YEAR FROM s.date) = (SELECT year FROM vars)) AS performer_services,
   	-- Услуг получено (терминов) пользователем 
	(SELECT count(*) FROM main.record r WHERE r.user_id = vars.uid 
	and r.done = true  AND EXTRACT(YEAR FROM r.date) = (SELECT year FROM vars)) AS recieved,
	-- Услуг создано пользователем
	(select count(*) from main.record r
		left join main.service s on s.uuid = r.service_id
		where s.performer_id = vars.uid  AND EXTRACT(YEAR FROM s.date) = (SELECT year FROM vars)
		
	) as providet,
	-- Услуг выполнено пользователем
	(select count(*) from main.record r
		left join main.service s on s.uuid = r.service_id
		where s.performer_id = vars.uid and r.done=true  AND EXTRACT(YEAR FROM r.date) = (SELECT year FROM vars)
		
	) as done_termins,
	-- Услуг подтверждено получателем
	(select count(*) from main.record r
		left join main.service s on s.uuid = r.service_id
		where s.performer_id = vars.uid and r.done=true 
		and r.user_confirm=true  AND EXTRACT(YEAR FROM r.date) = (SELECT year FROM vars)
		
	) as confirm_termins,
	-- Создано промоуслуг
	( select count(*) from main.service s
		where s.performer_id = vars.uid and s.promocode != ''  
		AND EXTRACT(YEAR FROM s.date) = (SELECT year FROM vars)
	) as created_promoservices,
	-- Комментариев в услугах пользователя
	(select count(*) from main.comment c
		left join main.record r on r.uuid = c.record_id
		left join main.service s on s.uuid = r.service_id
		where performer_id = vars.uid  AND EXTRACT(YEAR FROM (now() AT TIME ZONE 'UTC')) = (SELECT year FROM vars)
	) as comments_in_userservices
	
FROM vars;

`

	if err := r.db.QueryRow(query, user, year).Scan(
		&st.UserServices,
		&st.Recieved,
		&st.Provided,
		&st.Done,
		&st.Confirm,
		&st.Promoservises,
		&st.Comments,
	); err != nil {
		return st, err
	}
	return st, nil
}
