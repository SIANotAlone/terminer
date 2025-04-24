package repository

import (
	"fmt"
	"terminer/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

const (
	servicesTable       = "main.service"
	availableForTable   = "main.available_for"
	available_timeTable = "main.available_time"
	service_typeTable   = "main.service_type"
	time_layout         = "15:04"
)

type OfferingPostgres struct {
	db *sqlx.DB
}

func NewOfferingPostgres(db *sqlx.DB) *OfferingPostgres {
	return &OfferingPostgres{db: db}
}

func (r *OfferingPostgres) CreateOffering(offering models.NewService) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	var id uuid.UUID
	create_service_query := fmt.Sprintf("INSERT INTO %s (name, description, date, date_end, service_type_id, performer_id, available_for_all) VALUES ($1, $2, CURRENT_DATE, $3, $4, $5, $6) RETURNING uuid", servicesTable)
	row := tx.QueryRow(create_service_query, offering.Service.Name, offering.Service.Description, offering.Service.DateEnd, offering.Service.ServiceType, offering.Service.PerformerID, offering.Service.Available_for_all)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return id, err
	}
	for _, value := range offering.Available_for {
		create_avalable_for_query := fmt.Sprintf("INSERT INTO %s (service_id, user_id) VALUES ($1, $2)", availableForTable)
		_, err := tx.Exec(create_avalable_for_query, id, value.UserID)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
	}
	for _, value := range offering.Available_time {
		create_available_time_query := fmt.Sprintf("INSERT INTO %s (service_id, time_start, time_end) VALUES ($1, $2, $3)", available_timeTable)
		time_start, err := time.Parse(time_layout, value.TimeStart)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		time_end, err := time.Parse(time_layout, value.TimeEnd)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		_, err = tx.Exec(create_available_time_query, id, time_start.Format(time_layout), time_end.Format(time_layout))
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
	}

	return id, tx.Commit()
}

func (r *OfferingPostgres) UpdateService(service models.ServiceUpdate) error {

	query := fmt.Sprintf("UPDATE %s SET name = $1, description = $2, date_end = $3, service_type_id = $4 WHERE uuid = $5", servicesTable)
	row := r.db.QueryRow(query, service.Name, service.Description, service.DateEnd, service.ServiceType, service.UUID)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) DeleteService(id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE uuid = $1", servicesTable)
	row := r.db.QueryRow(query, id)
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) GetTypes() ([]models.ServiceType, error) {
	query := fmt.Sprintf("SELECT id, name FROM %s", service_typeTable)
	row, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var service_types []models.ServiceType
	for row.Next() {
		var service_type models.ServiceType
		if err := row.Scan(&service_type.ID, &service_type.Name); err != nil {
			return nil, err
		}
		service_types = append(service_types, service_type)
	}
	return service_types, nil
}

func (r *OfferingPostgres) GetServiceOwner(id uuid.UUID) (uuid.UUID, error) {
	query := fmt.Sprintf("SELECT performer_id FROM %s WHERE uuid = $1", servicesTable)
	row := r.db.QueryRow(query, id)
	var owner_id uuid.UUID
	if err := row.Scan(&owner_id); err != nil {
		return owner_id, err
	}
	return owner_id, nil
}

func (r *OfferingPostgres) CreateServiceType(s models.ServiceType) error {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1)", service_typeTable)
	_, err := r.db.Exec(query, s.Name)
	if err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) GetMyServices(user_id uuid.UUID) ([]models.MyService, error) {
	query := `select dc.uuid, dc.name, dc.description, dc.date, dc.date_end, st.name as service_type from main.service dc
	left join main.user u on dc.performer_id = u.uuid
	left join main.service_type st on dc.service_type_id = st.id 
	where dc.performer_id = $1 and date_end > CURRENT_DATE`
	var myservices []models.MyService
	row, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var ms models.MyService
		if err := row.Scan(&ms.ID, &ms.Name, &ms.Description, &ms.Date, &ms.DateEnd, &ms.ServiceType); err != nil {
			return nil, err
		}
		myservices = append(myservices, ms)
	}

	return myservices, nil
}

func (r *OfferingPostgres) GetAvailableService(user_id uuid.UUID) ([]models.AvailableService, error) {
	query := `select s.uuid as id,
 s.name as service, s.description, uu.first_name as p1, uu.last_name as p2,
 uu.email as p3, s.date, s.date_end, st.name as service_type
from main.available_for dc

left join main.user u on dc.user_id = u.uuid
left join main.service s on s.uuid = dc.service_id
left join main.user uu on s.performer_id = uu.uuid
left join main.service_type st on s.service_type_id = st.id

where dc.user_id = $1 and s.date_end >= CURRENT_DATE 
-- Фильтр по количеству доступных записей
and  (select count(*) from main.record r where r.service_id = s.uuid) < (select count(*) from main.available_time a where a.service_id =s.uuid )

union 
select s.uuid as id, s.name as service, s.description, u.first_name as p1,
u.last_name as p2, u.email as p3, s.date, s.date_end, st.name as service_type
from main.service s

left join main.user u on u.uuid = s.performer_id
left join main.service_type st on st.id = s.service_type_id

where s.available_for_all = true and date_end >= CURRENT_DATE 
and s.performer_id != $1
-- Фильтр по количеству доступных записей
and  (select count(*) from main.record r where r.service_id = s.uuid) < (select count(*) from main.available_time a where a.service_id =s.uuid )

`
	var available_services []models.AvailableService

	row, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var ms models.AvailableService
		if err := row.Scan(&ms.ID, &ms.Name, &ms.Description, &ms.FirstName, &ms.LastName, &ms.Email, &ms.Date, &ms.DateEnd, &ms.ServiceType); err != nil {
			return nil, err
		}
		available_services = append(available_services, ms)
	}
	return available_services, nil
}

func (r *OfferingPostgres) GetAvailableTime(service_id uuid.UUID) ([]models.ServiceAvailableTime, error) {
	query := `select id, service_id, time_start, time_end, booked 
		from main.available_time
		where service_id = $1 and booked = false`
	var available_times []models.ServiceAvailableTime
	row, err := r.db.Query(query, service_id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var available_time models.ServiceAvailableTime
		if err := row.Scan(&available_time.ID, &available_time.ServiceID, &available_time.TimeStart,
			&available_time.TimeEnd, &available_time.Booked); err != nil {
			return nil, err
		}
		available_times = append(available_times, available_time)
	}
	return available_times, nil
}

func (r *OfferingPostgres) CreatePromoService(offering models.NewPromoService) (uuid.UUID, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}
	var id uuid.UUID
	tx.Exec(`CREATE OR REPLACE FUNCTION generate_promo_code(length INT DEFAULT 8) RETURNS TEXT AS $$
DECLARE
    chars TEXT := 'ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    result TEXT := '';
    i INT;
BEGIN
    FOR i IN 1..length LOOP
        result := result || substr(chars, floor(random() * length(chars) + 1)::INT, 1);
    END LOOP;
    RETURN result;
END;
$$ LANGUAGE plpgsql;`)
	create_service_query := fmt.Sprintf(`INSERT INTO %s (name, description, date, date_end, service_type_id, performer_id,  promocode) VALUES ($1, $2, CURRENT_DATE, $3, $4, $5, generate_promo_code(10)) RETURNING uuid`, servicesTable)
	row := tx.QueryRow(create_service_query, offering.PromoService.Name, offering.PromoService.Description, offering.PromoService.DateEnd, offering.PromoService.ServiceType, offering.PromoService.PerformerID)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return id, err
	}

	for _, value := range offering.Available_time {
		create_available_time_query := fmt.Sprintf("INSERT INTO %s (service_id, time_start, time_end) VALUES ($1, $2, $3)", available_timeTable)
		time_start, err := time.Parse(time_layout, value.TimeStart)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		time_end, err := time.Parse(time_layout, value.TimeEnd)
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
		_, err = tx.Exec(create_available_time_query, id, time_start.Format(time_layout), time_end.Format(time_layout))
		if err != nil {
			tx.Rollback()
			return uuid.Nil, err
		}
	}

	return id, tx.Commit()
}

func (r *OfferingPostgres) GetUserTelegramID(user_id uuid.UUID) (string, error) {
	query := fmt.Sprintf("select telegram_chat_id from %s where uuid = $1", usersTable)
	row := r.db.QueryRow(query, user_id)
	var telegram_id string
	if err := row.Scan(&telegram_id); err != nil {
		return "", err
	}
	return telegram_id, nil
}

func (r *OfferingPostgres) GetAllUsersTelegramID() ([]string, error) {
	query := fmt.Sprintf("select telegram_chat_id from %s where telegram_chat_id IS NOT NULL", usersTable)

	row, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	var telegram_ids []string
	for row.Next() {
		var telegram_id string
		if err := row.Scan(&telegram_id); err != nil {
			return nil, err
		}
		telegram_ids = append(telegram_ids, telegram_id)
	}
	return telegram_ids, nil
}

func (r *OfferingPostgres) GetPromoCodeInfo(code string) (models.PromocodeInfo, error) {
	query := `select dc.uuid, dc.date, dc.date_end, dc.name, dc.description, max(u.last_name||' '||u.first_name) as performer, dc.promocode, count(f) as available_for
from main.service dc 
left join main.available_for f on f.service_id = dc.uuid
left join main.user u on u.uuid = dc.performer_id
where dc.promocode = $1 
group by 1`

	var info models.PromocodeInfo
	row := r.db.QueryRow(query, code)
	if err := row.Scan(&info.Service_ID, &info.Date, &info.Date_end, &info.Name, &info.Description, &info.Performer, &info.Promocode, &info.Available_for); err != nil {
		return info, err
	}
	return info, nil

}

func (r *OfferingPostgres) ActivatePromoCode(service_id uuid.UUID, user_id uuid.UUID) error {
	query := `insert into main.available_for(service_id, user_id) values ($1, $2)`

	_, err := r.db.Exec(query, service_id, user_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *OfferingPostgres) GetMyActualServices(user_id uuid.UUID) ([]models.MyActualService, error) {
	query := `SELECT 
    dc.uuid, 
    dc.name, 
    dc.description,
	st.name as service_type,
	dc.date, 
	dc.date_end,
	u.last_name || ' ' || u.first_name as performer,
    COALESCE(t.count_all, 0) AS count_all, 
    COALESCE(tt.count_free, 0) AS count_available --Доступні часи запису на послугу
FROM main.service dc
LEFT JOIN (
    SELECT service_id, COUNT(*) AS count_all 
    FROM main.available_time 
    GROUP BY service_id
) t ON t.service_id = dc.uuid
LEFT JOIN (
    SELECT service_id, COUNT(*) AS count_free 
    FROM main.available_time 
    WHERE booked = false
    GROUP BY service_id
) tt ON tt.service_id = dc.uuid
left join main.user u on u.uuid = dc.performer_id
left join main.service_type st on st.id =dc.service_type_id
WHERE dc.performer_id = $1 
 and COALESCE(tt.count_free, 0)>0 and dc.date_end < CURRENT_DATE

`

	var services []models.MyActualService
	row, err := r.db.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var service models.MyActualService
		if err := row.Scan(&service.ID, &service.Name, &service.Description, &service.ServiceType, &service.Date, &service.DateEnd, &service.Performer, &service.TotalSlots, &service.AvailableSlots); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *OfferingPostgres) GetHistoryMyServices(user_id uuid.UUID, limit int64, offset int64) ([]models.MyActualService, error) {
	query := `SELECT 
    dc.uuid, 
    dc.name, 
    dc.description,
	st.name as service_type,
	dc.date, 
	dc.date_end,
	u.last_name || ' ' || u.first_name as performer,
    COALESCE(t.count_all, 0) AS count_all, 
    COALESCE(tt.count_free, 0) AS count_available --Доступні часи запису на послугу
FROM main.service dc

LEFT JOIN (
    SELECT service_id, COUNT(*) AS count_all 
    FROM main.available_time 
    GROUP BY service_id
) t ON t.service_id = dc.uuid
LEFT JOIN (
    SELECT service_id, COUNT(*) AS count_free 
    FROM main.available_time 
    WHERE booked = false
    GROUP BY service_id
) tt ON tt.service_id = dc.uuid
left join main.user u on u.uuid = dc.performer_id
left join main.service_type st on st.id =dc.service_type_id
WHERE dc.performer_id = $1 
 LIMIT $2 OFFSET $3

`

	var services []models.MyActualService
	row, err := r.db.Query(query, user_id, limit, offset)
	if err != nil {
		return nil, err
	}
	for row.Next() {
		var service models.MyActualService
		if err := row.Scan(&service.ID, &service.Name, &service.Description, &service.ServiceType, &service.Date, &service.DateEnd, &service.Performer, &service.TotalSlots, &service.AvailableSlots); err != nil {
			return nil, err
		}
		services = append(services, service)
	}
	return services, nil

}
func (r *OfferingPostgres) GetTotalUserServices(user_id uuid.UUID) (int64, error) {
	query := `SELECT 
    count(*)
FROM main.service dc

WHERE dc.performer_id = $1
 `

	var count int64
	row := r.db.QueryRow(query, user_id)
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *OfferingPostgres) GetServicePromocode(service_id uuid.UUID) (models.PromocodeServiceInfo, error) {
	query := `select dc.uuid, dc.name, dc.description, dc.date, dc.date_end, st.name as servicetype, dc.promocode 
from main.service dc
left join main.service_type st on st.id = dc.service_type_id where dc.uuid = $1`

	var info models.PromocodeServiceInfo
	row := r.db.QueryRow(query, service_id)
	if err := row.Scan(&info.Service_id, &info.Name, &info.Description, &info.Date, &info.Date_end, &info.ServiceType, &info.Promocode); err != nil {
		return models.PromocodeServiceInfo{}, err
	}

	return info, nil
}
