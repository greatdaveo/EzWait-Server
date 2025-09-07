package models

import "time"

type Review struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	BookingID  uint      `gorm:"index" json:"booking_id"`
	Booking    Booking   `gorm:"foreignKey:BookingID;references:ID;constraint:OnDelete:CASCADE;" json:"booking"`
	CustomerID uint      `gorm:"index" json:"customer_id"`
	Customer   User      `gorm:"foreignKey:CustomerID;references:ID;constraint:OnDelete:CASCADE;" json:"customer"`
	StylistID  uint      `gorm:"index" json:"stylist_id"`
	Stylist    User      `gorm:"foreignKey:StylistID;references:ID;constraint:OnDelete:CASCADE;" json:"stylist"`
	Rating     int       `json:"rating"`
	Comment    string    `json:"comment"`
	CreatedAt  time.Time `json:"created_at"`
}
