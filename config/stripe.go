package config

import (
	"os"

	"github.com/stripe/stripe-go/v74"
)

func InitStripe() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	StripeWebhookSecret = os.Getenv("STRIPE_WEBHOOK_SECRET")
}
