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

func ViewStylistProfile(c *fiber.Ctx) error {
	stylistID := c.Params("stylistId")

	var stylist models.Stylist
	if err := config.DB.Where("stylist_id = ?", stylistID).First(&stylist).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Stylist not found",
		})
	}

	// To convert JSONB (Byte array) in PostgreSQL DB to Go structs
	var services []models.Service
	if err := json.Unmarshal(stylist.Services, &services); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse services",
		})
	}

	var serviceImages []string
	if err := json.Unmarshal(stylist.SampleOfServiceImg, &serviceImages); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse servicesImages",
		})
	}

	var timeSlots []string
	if err := json.Unmarshal(stylist.AvailableTimeSlots, &timeSlots); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse available time slot",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Stylist profile retrieved successfully",
		"data": fiber.Map{
			"id":                      stylist.ID,
			"stylist_id":              stylist.StylistID,
			"active_status":           stylist.ActiveStatus,
			"profile_picture":         stylist.ProfilePicture,
			"ratings":                 stylist.Ratings,
			"services":                services,
			"sample_of_service_img":   serviceImages,
			"available_time_slots":    timeSlots,
			"no_of_customer_bookings": stylist.NoOfCustomerBookings,
			"no_of_current_customers": stylist.NoOfCurrentCustomers,
			"created_at":              stylist.CreatedAt,
		},
	})
}

// PARTIAL UPDATE
func UpdateStylistProfile(c *fiber.Ctx) error {
	stylistIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "Invalid user ID format"})
	}
	stylistID := uint(stylistIDFloat)

	var input struct {
		ProfilePicture       *string           `json:"profile_picture"`
		Services             *[]models.Service `json:"services"`
		SampleOfServiceImg   *[]string         `json:"sample_of_service_img"`
		AvailableTimeSlots   *[]string         `json:"available_time_slots"`
		ActiveStatus         *bool             `json:"active_status"`
		NoOfCurrentCustomers *int              `json:"no_of_current_customers"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	// To know if stylist exist
	var stylist models.Stylist

	if err := config.DB.Where("stylist_id = ?", stylistID).First(&stylist).Error; err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Stylist not found" + err.Error(),
		})
	}

	// To convert JSONB to JSON
	if input.Services != nil {
		servicesJSON, _ := json.Marshal(input.Services)
		stylist.Services = servicesJSON
	}

	if input.SampleOfServiceImg != nil {
		imagesJSON, _ := json.Marshal(input.SampleOfServiceImg)
		stylist.SampleOfServiceImg = imagesJSON
	}

	if input.AvailableTimeSlots != nil {
		timeSlotsJSON, _ := json.Marshal(input.AvailableTimeSlots)
		stylist.AvailableTimeSlots = timeSlotsJSON
	}

	// To update only the provided fields
	if input.ProfilePicture != nil {
		stylist.ProfilePicture = *input.ProfilePicture
	}

	if input.ActiveStatus != nil {
		stylist.ActiveStatus = *input.ActiveStatus
	}

	if input.NoOfCurrentCustomers != nil {
		stylist.NoOfCurrentCustomers = *input.NoOfCurrentCustomers
	}

	// To save the updated Profile to DB
	if err := config.DB.Save(&stylist).Error; err != nil {

		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update profile" + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Profile updated successfully",
		"data":    stylist,
	})

}

// FULL UPDATE
func EditStylistProfile(c *fiber.Ctx) error {
	stylistIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	stylistID := uint(stylistIDFloat)
	// To parse the request body
	var input struct {
		ProfilePicture       string           `json:"profile_picture"`
		Services             []models.Service `json:"services"`
		SampleOfServiceImg   []string         `json:"sample_of_service_img"`
		AvailableTimeSlots   []string         `json:"available_time_slots"`
		ActiveStatus         bool             `json:"active_status"`
		NoOfCurrentCustomers int              `json:"no_of_current_customers"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	// To find stylist in database
	var stylist models.Stylist
	if err := config.DB.Where("stylist_id = ?", stylistID).First(&stylist).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Stylist not found",
		})
	}

	// To convert JSONB fields to JSON
	servicesJSON, _ := json.Marshal(input.Services)
	timeSlotsJSON, _ := json.Marshal(input.AvailableTimeSlots)
	imagesJSON, _ := json.Marshal(input.SampleOfServiceImg)

	// To overwrite all fields
	stylist.ProfilePicture = input.ProfilePicture
	stylist.Services = servicesJSON
	stylist.SampleOfServiceImg = imagesJSON
	stylist.AvailableTimeSlots = timeSlotsJSON
	stylist.ActiveStatus = input.ActiveStatus
	stylist.NoOfCurrentCustomers = input.NoOfCurrentCustomers

	// To save the updated profile
	if err := config.DB.Save(&stylist).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to edit profile",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Profile edited successfully",
		"data":    stylist,
	})
}

