package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

// Struct for services (array of objects)
type Service struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// Custom JSONB Type for pgxpool
type JSONB []byte

// Scan implements the sql.Scanner interface for JSONB.
func (j *JSONB) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	*j = JSONB(b) // Store the raw JSON bytes
	return nil
}

// Value implements the driver.Valuer interface for JSONB.
func (j JSONB) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return []byte(j), nil
}

// Stylist Model
type Stylist struct {
	ID                   uint      `json:"id"`
	StylistID            uint      `json:"stylist_id"`
	ActiveStatus         bool      `json:"active_status"`
	ProfilePicture       string    `json:"profile_picture"`
	Ratings              float64   `json:"ratings"`
	ServiceNameAndPrice  []Service `json:"services"` // Use []Service here
	SampleOfServiceImg   JSONB     `json:"sample_of_service_img"`
	AvailableTimeSlots   JSONB     `json:"available_time_slots"`
	NoOfCustomerBookings int       `json:"no_of_customer_bookings"`
	NoOfCurrentCustomers int       `json:"no_of_current_customers"`
	CreatedAt            time.Time `json:"created_at"`
}
