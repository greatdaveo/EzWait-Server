package middleware

import (
	"ezwait/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
)

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
