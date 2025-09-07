package services

import (
	"errors"
	"ezwait/config"
	"ezwait/internal/models"
	"fmt"
	"time"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
	"github.com/stripe/stripe-go/v74/refund"
)

type PaymentService struct{}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// To create payment intent on Stripe
func (ps *PaymentService) CreatePaymentIntent(amount int64, currency string, metadata map[string]string) (*stripe.PaymentIntent, error) {
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
	}

	// for key, value := range metadata {
	// 	params.AddMetadata(key, value)
	// }

	paymentIntent, err := paymentintent.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %v", err)
	}

	return paymentIntent, nil
}

// To confirm payment intent
func (ps *PaymentService) ConfirmPaymentIntent(paymentIntentID string) (*stripe.PaymentIntent, error) {
	paymentIntent, err := paymentintent.Confirm(paymentIntentID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to confirm payment intent: %v", err)
	}

	return paymentIntent, nil
}

// To process (50%) deposit payment
func (ps *PaymentService) ProcessDepositPayment(booking *models.Booking, customerEmail string) (*stripe.PaymentIntent, error) {
	if booking.DepositAmount == nil || *booking.DepositAmount <= 0 {
		return nil, errors.New("invalid deposit amount")
	}

	finalAmount := *booking.TotalAmount - *booking.DepositAmount
	if finalAmount < 0 {
		return nil, errors.New("invalid final payment amount")
	}

	// To convert to pence (stripe smallest currency unit)
	amountInPence := int64(*booking.DepositAmount * 100)

	metadata := map[string]string{
		"booking_id":     fmt.Sprintf("%d", booking.ID),
		"payment_type":   "deposit",
		"customer_id":    fmt.Sprintf("%d", booking.UserID),
		"stylist_id":     fmt.Sprintf("%d", booking.StylistID),
		"customer_email": customerEmail,
	}

	paymentIntent, err := ps.CreatePaymentIntent(amountInPence, "gbp", metadata)
	if err != nil {
		return nil, err
	}

	// To create payment record
	payment := models.Payment{
		BookingID:             booking.ID,
		Amount:                *booking.DepositAmount,
		PaymentType:           "deposit",
		PaymentStatus:         "pending",
		StripePaymentIntentID: &paymentIntent.ID,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment record %v", err)
	}

	return paymentIntent, nil
}

// To process the 50% deposit payment
func (ps *PaymentService) ProcessFinalPayment(booking *models.Booking, customerEmail string) (*stripe.PaymentIntent, error) {
	if booking.TotalAmount == nil || booking.DepositAmount == nil {
		return nil, errors.New("invalid payment amounts")
	}

	finalAmount := *booking.TotalAmount - *booking.DepositAmount
	if finalAmount <= 0 {
		return nil, errors.New("invalid final payment amount")
	}

	// To convert to pence
	amountInPence := int64(finalAmount * 100)

	metadata := map[string]string{
		"booking_id":     fmt.Sprintf("%d", booking.ID),
		"payment_type":   "final",
		"customer_id":    fmt.Sprintf("%d", booking.UserID),
		"stylist_id":     fmt.Sprintf("%d", booking.StylistID),
		"customer_email": customerEmail,
	}

	paymentIntent, err := ps.CreatePaymentIntent(amountInPence, "gbp", metadata)
	if err != nil {
		return nil, err
	}

	// To create payment record
	payment := models.Payment{
		BookingID:             booking.ID,
		Amount:                finalAmount,
		PaymentType:           "final",
		PaymentStatus:         "pending",
		StripePaymentIntentID: &paymentIntent.ID,
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		return nil, fmt.Errorf("failed to create payment record: %v", err)
	}

	return paymentIntent, nil
}

// To process refund
func (ps *PaymentService) ProcessRefund(paymentIntentID string, amount int64, reason string) (*stripe.Refund, error) {
	params := &stripe.RefundParams{
		PaymentIntent: stripe.String(paymentIntentID),
		Amount:        stripe.Int64(amount),
		Reason:        stripe.String(reason),
	}

	refund, err := refund.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to process refund: %v", err)
	}

	return refund, nil
}

// To update payment status
func (ps *PaymentService) UpdatePaymentStatus(paymentIntentID string, status string) error {
	var payment models.Payment
	if err := config.DB.Where("stripe_payment_intent_id = ?", paymentIntentID).First(&payment).Error; err != nil {
		return fmt.Errorf("payment not found: %v", err)
	}

	payment.PaymentStatus = status
	payment.UpdatedAt = time.Now()

	if err := config.DB.Save(&payment).Error; err != nil {
		return fmt.Errorf("failed to update payment status: %v", err)
	}

	// To update the booking status
	var booking models.Booking
	if err := config.DB.First(&booking, payment.BookingID).Error; err != nil {
		return fmt.Errorf("booking not found: %v", err)
	}

	if payment.PaymentType == "deposit" && status == "completed" {
		booking.DepositPaid = true
	} else if payment.PaymentType == "final" && status == "completed" {
		booking.FinalPaymentPaid = true
	}

	if err := config.DB.Save(&booking).Error; err != nil {
		return fmt.Errorf("failed to update booking: %v", err)
	}

	return nil
}
