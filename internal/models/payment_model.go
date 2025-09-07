package models

import "time"

type Payment struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	BookingID             uint      `gorm:"index" json:"booking_id"`
	Booking               Booking   `gorm:"foreignKey:BookingID;references:ID;constraint:OnDelete:CASCADE;" json:"booking"`
	Amount                float64   `json:"amount"`
	PaymentType           string    `json:"payment_type"`
	PaymentStatus         string    `json:"payment_status"`
	StripePaymentIntentID *string   `json:"stripe_payment_intent_id"`
	StripeRefundID        *string   `json:"stripe_refund_id"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
