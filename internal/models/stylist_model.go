package models

import (
	"encoding/json"
	"time"
)

type Service struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type SampleOfService struct {
	Img     string `json:"img_url"`
	Caption string `json:"caption"`
}

type Stylist struct {
	ID                   uint            `gorm:"primaryKey" json:"id"`
	StylistID            uint            `gorm:"uniqueIndex" json:"stylist_id"`
	User                 User            `gorm:"foreignKey:StylistID;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
	ActiveStatus         bool            `json:"active_status"`
	ProfilePicture       string          `json:"profile_picture"`
	Ratings              float64         `json:"ratings"`
	Services             json.RawMessage `json:"services" gorm:"type:jsonb"`
	SampleOfServices     json.RawMessage `json:"sample_of_services" gorm:"type:jsonb"`
	AvailableTimeSlots   json.RawMessage `json:"available_time_slots" gorm:"type:jsonb"`
	NoOfCustomerBookings int             `json:"no_of_customer_bookings"`
	NoOfCurrentCustomers int             `json:"no_of_current_customers"`
	AutoConfirm          bool            `json:"auto_confirm" gorm:"default:false"`
	CreatedAt            time.Time       `json:"created_at"`
}

