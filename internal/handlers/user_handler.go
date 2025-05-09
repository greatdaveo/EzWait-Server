package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"

	"github.com/gofiber/fiber/v2"
)

func UpdateUserProfile(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid user ID format"})
	}
	userID := uint(userIDFloat)

	var input struct {
		Name           string `json:"name"`
		Email          string `json:"email"`
		Number         string `json:"number"`
		Location       string `json:"location"`
		ProfilePicture string `json:"profile_picture"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Number = input.Number
	user.Location = input.Location
	user.ProfilePicture = input.ProfilePicture

	if err := config.DB.Save(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update user profile",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"data":    user,
	})
}
