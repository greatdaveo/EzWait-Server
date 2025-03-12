package models

import "time"

type Booking struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	StylistID     uint      `json:"stylist_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	BookingDay    time.Time `json:"booking_day"`
	BookingStatus string    `json:"booking_status"`
	CreatedAt     time.Time `json:"created_at"`
}
