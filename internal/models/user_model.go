package models

import "time"

const (
	RoleStylist  = "stylist"
	RoleCustomer = "customer"
)

type User struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Number          int       `json:"number"`
	Role            string    `json:"role"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	Location        string    `json:"location"`
	CreatedAt       time.Time `json:"created_at"`
}
