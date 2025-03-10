package models

import (
	"gorm.io/gorm"
)

const (
	RoleStylist  = "stylist"
	RoleCustomer = "customer"
)

// User model
type User struct {
	gorm.Model
	Name            string `json:"name"`
	Email           string `json:"email" gorm:"uniqueIndex"`
	Number          string `json:"number"`
	Role            string `json:"role"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password" gorm:"-"`
	Location        string `json:"location"`
}
