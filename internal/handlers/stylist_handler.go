package handlers

import (
	"encoding/json"
	"ezwait/config"
	"ezwait/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateStylistProfile(c *fiber.Ctx) error {

	stylistIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid user ID format"})
	}
	stylistID := uint(stylistIDFloat)

	var existingStylist models.Stylist
	if err := config.DB.Where("stylist_id = ?", stylistID).First(&existingStylist).Error; err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Stylist profile already exists"})
	}

	// To extract data from the JSON req
	var input struct {
		ProfilePicture     string           `json:"profile_picture"`
		Services           []models.Service `json:"services"`
		SampleOfServiceImg []string         `json:"sample_of_service_img"`
		AvailableTimeSlots []string         `json:"available_time_slots"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input: " + err.Error()})
	}

	// To Convert slices to JSON
	servicesJSON, err := json.Marshal(input.Services)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode services"})
	}

	imagesJSON, err := json.Marshal(input.SampleOfServiceImg)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode images"})
	}

	timeSlotsJSON, err := json.Marshal(input.AvailableTimeSlots)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to encode time slots"})
	}

	// To create Stylist Profile
	stylist := models.Stylist{
		StylistID:            stylistID,
		ActiveStatus:         true,
		ProfilePicture:       input.ProfilePicture,
		Ratings:              0.0,
		Services:             servicesJSON,
		SampleOfServiceImg:   imagesJSON,
		AvailableTimeSlots:   timeSlotsJSON,
		NoOfCustomerBookings: 0,
		NoOfCurrentCustomers: 0,
		CreatedAt:            time.Now(),
	}

	// To save to database
	if err := config.DB.Create(&stylist).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create stylist profile: " + err.Error()})
	}
	// Success Response
	return c.Status(201).JSON(fiber.Map{
		"message": "Stylist profile created successfully",
		"data":    stylist,
	})
}

