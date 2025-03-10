package handlers

import (
	"encoding/json"
	"ezwait/config"
	"ezwait/internal/models"

	"github.com/gofiber/fiber/v2"
)

func CreateStylistProfile(c *fiber.Ctx) error {
	// Get the stylist ID from middleware
	stylistID, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(401).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Check if stylist profile already exists
	var exists bool
	err := config.DB.QueryRow(
		c.Context(),
		"SELECT EXISTS (SELECT 1 FROM stylists WHERE stylist_id=$1)", stylistID,
	).Scan(&exists)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to check stylist profile"})
	}

	if exists {
		return c.Status(400).JSON(fiber.Map{"error": "Stylist profile already exists"})
	}

	// Parse request body
	var input struct {
		ProfilePicture     string           `json:"profile_picture"`
		Services           []models.Service `json:"services"` // Use []models.Service here
		SampleOfServiceImg []string         `json:"sample_of_service_img"`
		AvailableTimeSlots []string         `json:"available_time_slots"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input: " + err.Error()})
	}

	// Convert data to JSONB
	serviceData, _ := json.Marshal(input.Services)
	imageData, _ := json.Marshal(input.SampleOfServiceImg)
	timeSlotData, _ := json.Marshal(input.AvailableTimeSlots)

	// Insert into PostgreSQL
	_, err = config.DB.Exec(
		c.Context(),
		`INSERT INTO stylists 
		(stylist_id, active_status, profile_picture, ratings, services, service_img, available_time_slots, no_of_customer_bookings, no_of_current_customers, created_at) 
		VALUES ($1, $2, $3, $4, $5::jsonb, $6::jsonb, $7::jsonb, $8, $9, NOW())`,
		stylistID, true, input.ProfilePicture, 0.0, serviceData, imageData, timeSlotData, 0, 0,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create stylist profile: " + err.Error()})
	}

	// Return success
	return c.Status(201).JSON(fiber.Map{
		"message": "Stylist profile created successfully",
	})
}
