package models

const (
	RoleStylist  = "stylist"
	RoleCustomer = "customer"
)

// User model
type User struct {
	ID              uint   `gorm:"primaryKey" json:"id"`
	Name            string `json:"name"`
	Email           string `json:"email" gorm:"uniqueIndex"`
	Number          string `json:"number"`
	Role            string `json:"role"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password" gorm:"-"`
	Location        string `json:"location"`
	Stylist         *Stylist `gorm:"foreignKey:StylistID;references:ID"`
}
