package models

type Stylist struct {
	ID                   uint    `json:"id"`
	UserID               uint    `json:"user_id"`
	Ratings              float64 `json:"ratings"`
	Services             string  `json:"services"`
	Availability         string  `json:"availability"`
	NoOfCustomerBookings int     `json:"no_of_customer_bookings"`
	NoOfCurrentCustomers int     `json:"no_of_current_customers"`
	CreatedAt            string  `json:"created_at"`
}
