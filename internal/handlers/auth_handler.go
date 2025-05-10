package handlers

import (
	"errors"
	"ezwait/config"
	"ezwait/internal/models"
	"ezwait/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var store = session.New()

func RegisterHandler(c *fiber.Ctx) error {
	var user models.User

	// To parse request
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// To check if email exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already registered"})
	}

	// To validate role
	if user.Role != models.RoleStylist && user.Role != models.RoleCustomer {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid role"})
	}

	// To validate passwords
	if user.Password != user.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	// To hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	// To save user
	if err := config.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to register user",
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "User created successfully",
		"data": fiber.Map{
			"id":              user.ID,
			"name":            user.Name,
			"email":           user.Email,
			"number":          user.Number,
			"role":            user.Role,
			"location":        user.Location,
			"profile_picture": user.ProfilePicture,
		},
	})
}

func LoginHandler(c *fiber.Ctx) error {
	// Login form request structure
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest

	// To parse request
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// To find user by email
	var user models.User
	err := config.DB.Where("email = ?", loginReq.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Database error"})
	}

	// To compare hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid Credentials"})
	}

	// To generate JWT token
	token, err := utils.GenerateToken(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// Return response
	return c.Status(200).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"data": fiber.Map{
			"id":              user.ID,
			"name":            user.Name,
			"email":           user.Email,
			"number":          user.Number,
			"role":            user.Role,
			"location":        user.Location,
			"profile_picture": user.ProfilePicture,
		},
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	// To retrieve the session
	session, err := store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	// To destroy the session
	if err := session.Destroy(); err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}

func ChangePassword(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized user",
		})
	}

	userID := uint(userIDFloat)

	var input struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	if input.NewPassword != input.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{
			"error": "New passwords do not match",
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return c.Status(401).JSON(fiber.Map{
			"error": "Current password is incorrect",
		})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to hash new password",
		})
	}

	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update password",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Password updated successfully",
	})
}

func DeleteAccount(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	userID := uint(userIDFloat)

	var user models.User

	if err := config.DB.Delete(&user, userID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete account: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Your account has been deleted successfully",
	})
}
