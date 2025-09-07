package models

import "time"

type Booking struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"index" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE;" json:"user"`
	StylistID        uint      `gorm:"index;not null" json:"stylist_id"`
	Stylist          Stylist   `gorm:"foreignKey:StylistID;references:StylistID;constraint:OnDelete:CASCADE;" json:"stylist"`
	StartTime        time.Time `json:"start_time"`
	EndTime          time.Time `json:"end_time"`
	BookingDay       time.Time `json:"booking_day"`
	BookingStatus    string    `json:"booking_status"`
	TotalAmount      *float64  `json:"total_amount"`
	DepositAmount    *float64  `json:"deposit_amount"`
	DepositPaid      bool      `json:"deposit_paid" gorm:"default:false"`
	FinalPaymentPaid bool      `json:"final_payment_paid" gorm:"default:false"`
	ServiceName      string    `json:"service_name"`
	Notes            string    `json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
}
