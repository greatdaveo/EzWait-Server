package models

import "time"

type Booking struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
	StylistID     uint      `gorm:"index;not null" json:"stylist_id"`
	Stylist       Stylist   `gorm:"foreignKey:StylistID;references:ID;constraint:OnDelete:CASCADE;" json:"stylist"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	BookingDay    time.Time `json:"booking_day"`
	BookingStatus string    `json:"booking_status"`
	CreatedAt     time.Time `json:"created_at"`
}


