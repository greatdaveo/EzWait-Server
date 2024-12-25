package models

type Booking struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	StylistID     uint   `json:"stylist_id"`
	BookingTime   string `json:"booking_time"`
	BookingDay    string `json:"booking_day"`
	BookingStatus string `json:"booking_status"`
	CreatedAt     string `json:"created_at"`
}
