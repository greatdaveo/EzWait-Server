package handlers

import (
	"encoding/json"
	"ezwait/config"
	"ezwait/internal/models"
	"ezwait/internal/services"
	"io"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/webhook"
)

type PaymentHandler struct {
	paymentService *services.PaymentService
}

func NewPaymentHandler() *PaymentHandler {
	return &PaymentHandler{
		paymentService: services.NewPaymentService(),
	}
}

// To create deposit payment
func (ph *PaymentHandler) CreateDepositPayment(c *fiber.Ctx) error {
	userIdFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": "Unauthorized",
		})
	}

	userID := uint(userIdFloat)

	var input struct {
		BookingID uint `json:"booking_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input",
		})
	}

	// To get bookings details
	var booking models.Booking
	if err := config.DB.Preload("Stylist").First(&booking, input.BookingID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Booking not found",
		})
	}

	// To verify booking belongs to the user
	if booking.UserID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error":   true,
			"message": "Not authorized to pay for this booking",
		})
	}

	// To check if the deposit is already paid
	if booking.DepositPaid {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Deposit already paid",
		})
	}

	// To get user email
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	// To create payment intent
	paymentIntent, err := ph.paymentService.ProcessDepositPayment(&booking, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create payment intent",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Payment intent created",
		"data": fiber.Map{
			"client_secret":     paymentIntent.ClientSecret,
			"payment_intent_id": paymentIntent.ID,
			"amount":            paymentIntent.Amount,
			"currency":          paymentIntent.Currency,
		},
	})
}

// Create Final Payment
func (ph *PaymentHandler) CreateFinalPayment(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": "Unauthorized",
		})
	}
	userID := uint(userIDFloat)

	var input struct {
		BookingID uint `json:"booking_id"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid input",
		})
	}

	// To get booking details
	var booking models.Booking
	if err := config.DB.Preload("Stylist").First(&booking, input.BookingID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "Booking not found",
		})
	}

	// To verify booking  to user
	if booking.UserID != userID {
		return c.Status(403).JSON(fiber.Map{
			"error":   true,
			"message": "Not authorized to pay for this booking",
		})
	}

	// To check if deposit is paid
	if !booking.DepositPaid {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Deposit must be paid first",
		})
	}

	// To check if final payment are paid
	if booking.FinalPaymentPaid {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Final payment already paid",
		})
	}

	// To get user email
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error":   true,
			"message": "User not found",
		})
	}

	// To create payment intent
	paymentIntent, err := ph.paymentService.ProcessFinalPayment(&booking, user.Email)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to create payment intent",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Payment intent created",
		"data": fiber.Map{
			"client_secret":     paymentIntent.ClientSecret,
			"payment_intent_id": paymentIntent.ID,
			"amount":            paymentIntent.Amount,
			"currency":          paymentIntent.Currency,
		},
	})
}

// To handle Stripe webhook events
func (ph *PaymentHandler) StripeWebhook(c *fiber.Ctx) error {
	body, err := io.ReadAll(c.Request().BodyStream())
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to read request body",
		})
	}

	// To verify webhook signature
	event, err := webhook.ConstructEvent(body, c.Get("Stripe-Signature"), config.StripeWebhookSecret)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid webhook signature",
		})
	}

	// To handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to parse payment intent",
			})
		}

		// To update payment status
		if err := ph.paymentService.UpdatePaymentStatus(paymentIntent.ID, "completed"); err != nil {
			log.Printf("Failed to update payment status: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to update payment status",
			})
		}

	case "payment_intent.payment_failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to parse payment intent",
			})
		}

		// To update payment status
		if err := ph.paymentService.UpdatePaymentStatus(paymentIntent.ID, "failed"); err != nil {
			log.Printf("Failed to update payment status: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to update payment status",
			})
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Webhook processed successfully",
	})
}

// To return payment history for a user
func (ph *PaymentHandler) GetPaymentHistory(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error":   true,
			"message": "Unauthorized",
		})
	}
	userID := uint(userIDFloat)

	var payments []models.Payment
	if err := config.DB.Preload("Booking").
		Joins("JOIN bookings ON payments.booking_id = bookings.id").
		Where("bookings.user_id = ?", userID).
		Order("payments.created_at DESC").
		Find(&payments).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to fetch payment history",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error":   false,
		"message": "Payment history retrieved",
		"data":    payments,
	})
}
