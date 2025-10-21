package models

import (
	"time"

	"github.com/google/uuid"
)

type Service struct {
	UUID              uuid.UUID `json:"-"`
	Name              string    `json:"name" binding:"required" omitempty:"true"`
	Description       string    `json:"description" binding:"required" omitempty:"true"`
	Date              time.Time `json:"-"`
	DateEnd           time.Time `json:"date_end" binding:"required" omitempty:"true"`
	ServiceType       int       `json:"service_type" binding:"required" omitempty:"true"`
	PerformerID       uuid.UUID `json:"-"`
	Available_for_all bool      `json:"for_all" omitempty:"true"`
	MassageType       *int      `json:"massage_type"`
}

type Available_time struct {
	ID        int       `json:"-"`
	ServiceID uuid.UUID `json:"-`
	TimeStart string    `json:"time_start" binding:"required" omitempty:"true"`
	TimeEnd   string    `json:"time_end" binding:"required" omitempty:"true"`
}
type Available_for struct {
	ID        int       `json:"-"`
	ServiceID uuid.UUID `json:"-"`
	UserID    uuid.UUID `json:"user_id" binding:"required" omitempty:"true"`
}

type NewService struct {
	Service        Service          `json:"service" binding:"required" omitempty:"true"`
	Available_time []Available_time `json:"available_time"`
	Available_for  []Available_for  `json:"available_for"`
}
type PromoService struct {
	UUID        uuid.UUID `json:"-"`
	Name        string    `json:"name" binding:"required" omitempty:"true"`
	Description string    `json:"description" binding:"required" omitempty:"true"`
	Date        time.Time `json:"-"`
	DateEnd     time.Time `json:"date_end" binding:"required" omitempty:"true"`
	ServiceType int       `json:"service_type" binding:"required" omitempty:"true"`
	PerformerID uuid.UUID `json:"-"`
}
type NewPromoService struct {
	PromoService   PromoService     `json:"promoservice" binding:"required" omitempty:"true"`
	Available_time []Available_time `json:"available_time"`
}

type ServiceUpdate struct {
	UUID        uuid.UUID `json:"id" binding:"required" omitempty:"true"`
	Name        string    `json:"name" binding:"required" omitempty:"true"`
	Description string    `json:"description" binding:"required" omitempty:"true"`
	DateEnd     time.Time `json:"date_end" binding:"required" omitempty:"true"`
	ServiceType int       `json:"service_type" binding:"required" omitempty:"true"`
}

type ServiceDelete struct {
	UUID uuid.UUID `json:"id" binding:"required" omitempty:"true"`
}

type ServiceType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type MassageType struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	CasualName string `json:"casual_name"`
}

type MyService struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	DateEnd     time.Time `json:"date_end"`
	ServiceType string    `json:"service_type"`
}
type AvailableService struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	Date        time.Time `json:"date"`
	DateEnd     time.Time `json:"date_end"`
	ServiceType string    `json:"service_type"`
	MassageType string    `json:"massage_type"`
}

type ServiceAvailableTime struct {
	ID        int       `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`
	TimeStart time.Time `json:"time_start"`
	TimeEnd   time.Time `json:"time_end"`
	Booked    bool      `json:"booked"`
}

type ServiceAvailableTimeInput struct {
	ID uuid.UUID `json:"service_id" binding:"required" omitempty:"true"`
}

type PromocodeInfo struct {
	Service_ID    uuid.UUID `json:"service_id"`
	Date          time.Time `json:"date"`
	Date_end      time.Time `json:"date_end"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Performer     string    `json:"performer"`
	Promocode     string    `json:"-"`
	Available_for int64     `json:"-"`
}

type PromocodeValidation struct {
	Valid        bool          `json:"valid"`
	PromeService PromocodeInfo `json:"promoservice"`
}

type PromocodeValidationInput struct {
	Promocode string `json:"promocode" binding:"required" omitempty:"true"`
}

type PromocodeActivation struct {
	Promocode string    `json:"promocode" binding:"required" omitempty:"true"`
	User_ID   uuid.UUID `json:"-"`
}

type PromocodeActivationInput struct {
	Promocode string `json:"promocode" binding:"required" omitempty:"true"`
}

type MyActualService struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ServiceType string    `json:"service_type"`
	Date        time.Time `json:"date"`
	DateEnd     time.Time `json:"date_end"`
	Performer   string    `json:"performer"`
	TotalSlots  int64     `json:"total_slots"`
	BookedSlots int64     `json:"booked_slots"`
	MassageType string    `json:"massage_type"`
}

type MyHistoryServiceInput struct {
	Limit  int64 `json:"limit" omitempty:"true"`
	Offset int64 `json:"offset" omitempty:"true"`
}

type PromocodeServiceInfo struct {
	Service_id  uuid.UUID `json:"service_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Date_end    time.Time `json:"date_end"`
	ServiceType string    `json:"service_type"`
	Promocode   string    `json:"promocode"`
}

type UserServiceHistory struct {
	History []MyActualService `json:"history"`
	Total   int64             `json:"total"`
}

type FullServiceInformation struct {
	ServiceInformation  ServiceInformation    `json:"service"`
	Available_for       []Available_for_Info  `json:"available_for"`
	Available_time_Info []Available_time_Info `json:"available_time"`
}

type ServiceInformation struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Date            *time.Time `json:"date"`
	DateEnd         time.Time  `json:"date_end"`
	ServiceTypeID   int        `json:"service_type_id"`
	AvailableForAll bool       `json:"available_for_all"`
	MassageTypeID   int        `json:"massage_type_id"`
}

type Available_time_Info struct {
	ID        int    `json:"id"`
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Booked    bool   `json:"booked"`
}
type Available_for_Info struct {
	ServiceID uuid.UUID `json:"-"`
	ID        int       `json:"id"`
	UserID    uuid.UUID `json:"user_id" binding:"required" omitempty:"true"`
	Name      string    `json:"name"`
}
type NewAvailableTime struct {
	ServiceID uuid.UUID `json:"service_id" binding:"required" omitempty:"true"`
	TimeStart string    `json:"time_start" binding:"required" omitempty:"true"`
	TimeEnd   string    `json:"time_end" binding:"required" omitempty:"true"`
}

type DeleteAvailableTime struct {
	ID int `json:"id" binding:"required" omitempty:"true"`
}

type NewAvailableFor struct {
	UserID    uuid.UUID `json:"user_id" binding:"required" omitempty:"true"`
	ServiceID uuid.UUID `json:"service_id" binding:"required" omitempty:"true"`
}

type DeleteAvailableFor struct {
	ID int `json:"id" binding:"required" omitempty:"true"`
}
