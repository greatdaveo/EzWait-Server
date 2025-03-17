package handlers

import (
	"encoding/json"
	"ezwait/config"
	"ezwait/internal/models"
	"fmt"
	"strconv"
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
		StylistID          uint             `json:"stylist_id"`
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

	// To fetch the stylist details and the Stylist user data
	if err := config.DB.Preload("User").Where("stylist_id = ?", stylistID).First(&stylist).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stylist profile: " + err.Error()})
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

	// To fetch the associated user details of the stylist
	var user models.User
	if err := config.DB.Where("id = ?", stylist.StylistID).First(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch stylist profile: " + err.Error()})
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
		"user": fiber.Map{
			"id":       user.ID,
			"name":     user.Name,
			"email":    user.Email,
			"number":   user.Number,
			"role":     user.Role,
			"location": user.Location,
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

func ViewAllStylists(c *fiber.Ctx) error {
	// To extract optional query params
	serviceFilter := c.Query("service")
	sortBy := c.Query("sort")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("page", "10"))

	// For valid pagination values
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit < 50 {
		limit = 10
	}

	// To query the stylist table
	query := config.DB.Model(&models.Stylist{})

	if serviceFilter != "" {
		query = query.Where("services @> ?", fmt.Sprintf(`"%s"`, serviceFilter))
	}

	if sortBy == "ratings" {
		query = query.Order("ratings DESC")
	} else if sortBy == "name" {
		query = query.Order("name ASC")
	}

	offset := (page - 1) * limit
	query = query.Offset(offset).Limit(limit)

	// To fetch the queries from DB
	var stylists []models.Stylist
	if err := query.Find(&stylists).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch stylist: " + err.Error(),
		})
	}

	// To prepare response with additional user details
	var stylistResponse []map[string]interface{}

	for _, stylist := range stylists {
		var user models.User

		// To fetch the stylists name & location
		if err := config.DB.Where("id = ?", stylist.StylistID).First(&user).Error; err != nil {
			fmt.Println("Error fetching user details: ", err)
			continue
		}

		// To convert JSONB response to Go struct
		var services []models.Service
		var sampleImgs []string
		var timeSlots []string

		_ = json.Unmarshal(stylist.Services, &services)
		_ = json.Unmarshal(stylist.SampleOfServiceImg, &sampleImgs)
		_ = json.Unmarshal(stylist.AvailableTimeSlots, &timeSlots)

		// To append stylist info to response array
		stylistResponse = append(stylistResponse, map[string]interface{}{
			"id":                      stylist.ID,
			"stylist_id":              stylist.StylistID,
			"name":                    user.Name,
			"location":                user.Location,
			"active_status":           stylist.ActiveStatus,
			"profile_picture":         stylist.ProfilePicture,
			"ratings":                 stylist.Ratings,
			"services":                services,
			"sample_of_service_img":   sampleImgs,
			"available_time_slots":    timeSlots,
			"no_of_customer_bookings": stylist.NoOfCustomerBookings,
			"auto_confirm":            stylist.AutoConfirm,
			"created_at":              stylist.CreatedAt,
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Stylists retrieved successfully",
		"data":    stylistResponse,
		"page":    page,
		"limit":   limit,
	})

}
