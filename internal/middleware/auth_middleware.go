package middleware

import (
	"ezwait/internal/models"
	"ezwait/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// To check if the user is authenticated
func AuthMiddleware(c *fiber.Ctx) error {
	// Get token from Authorization header
	token := c.Get("Authorization")
	// To check if the token is provided
	if token == "" {
		return c.Status(401).JSON(fiber.Map{
			"error": "Missing token",
		})
	}

	// To extract the token
	token = strings.Replace(token, "Bearer ", "", -1)

	// To validate JWT
	claims, err := utils.VerifyToken(token)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	// To extract user and role from claims
	user, ok := (*claims)["user"].(float64)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user data - user id not found",
		})
	}

	role, ok := (*claims)["role"].(string)

	if !ok {
		// role = ""
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user data - role not found",
		})
	}

	// To add claims to locals for use in handlers
	c.Locals("user", user)

	c.Locals("role", role)

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
