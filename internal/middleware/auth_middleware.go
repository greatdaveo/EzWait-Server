package middleware

import (
	"ezwait/internal/models"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

// To set session store
func SetSessionStore(s *session.Store) {
	store = s
}

func AuthMiddleware(c *fiber.Ctx) error {
	// To get the session from the store
	session, err := store.Get(c)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// To check if the user exists in the session
	user := session.Get("user")
	fmt.Println(user)

	if user != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized, please login",
		})
	}

	// To pass the user and role to locals (req) for further use
	c.Locals("user", user)
	c.Locals("role", session.Get("role"))

	return c.Next()
}

func ValidateUser(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// To check for any missing fields
	if user.Name == "" || user.Email == "" || user.Password == "" || user.ConfirmPassword == "" || user.Role == "" {
		return c.Status(400).JSON(fiber.Map{"error": "All fields are required"})
	}

	// To validate email format
	if !strings.Contains(user.Email, "@") {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid email format"})
	}

	return c.Next()
}

// To ensure the user is a customer
func ValidateCustomer(c *fiber.Ctx) error {
	role := c.Locals("role")

	if role != "customer" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Only customers are permitted",
		})
	}

	return c.Next()
}

// To ensure the user is a customer
func ValidateStylist(c *fiber.Ctx) error {
	role := c.Locals("role")

	if role != "stylist" {
		return c.Status(403).JSON(fiber.Map{
			"error": "Only stylists are permitted",
		})
	}

	return c.Next()
}
