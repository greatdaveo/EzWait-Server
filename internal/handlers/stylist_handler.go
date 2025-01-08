package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func CreateStylistProfile(c *fiber.Ctx) error {
	// To get the stylistId from the local storage
	stylistID := c.Locals("user").(float64)

	// To check if the stylist profile already exist
	var existingStylist models.Stylist
	err := config.DB.QueryRow(
		c.Context(),
		"SELECT * FROM stylists WHERE stylist_id=$1",
		stylistID,
	).Scan(
		&existingStylist.ID,
		&existingStylist.StylistID,
		&existingStylist.Ratings,
		&existingStylist.Services,
		&existingStylist.Availability,
		&existingStylist.NoOfCustomerBookings,
		&existingStylist.NoOfCurrentCustomers,
		&existingStylist.CreatedAt,
	)

	if err == nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Stylist profile already exists: " + err.Error(),
		})
	}

	// To dynamically calculate the number of active customers bookings
	var noOfCustomerBookings int

	err = config.DB.QueryRow(
		c.Context(),
		`SELECT COUNT(*) FROM bookings
		 WHERE stylist_id=$1 AND booking_status NOT IN ("completed", "cancelled")`,
		stylistID,
	).Scan(&noOfCustomerBookings)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to calculate customer bookings: " + err.Error(),
		})
	}

	// To parse the req body
	var input struct {
		Services     string `json:"services"`
		Availability string `json:"availability"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	// To insert the stylist profile data
	_, err = config.DB.Exec(
		c.Context(),
		`INSERT INTO stylists 
		(stylist_id, ratings, services, availability, no_of_customer_bookings, no_of_current_customers, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, NOW())`,
		stylistID, 0.0, input.Services, input.Availability, noOfCustomerBookings, 0,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create stylist profile: " + err.Error(),
		})
	}

	// To fetch the created stylist profile
	var stylist models.Stylist
	err = config.DB.QueryRow(
		c.Context(),
		"SELECT * FROM stylists WHERE stylist_id=$1",
		stylistID,
	).Scan(
		&stylist.ID,
		&stylist.StylistID,
		&stylist.Ratings,
		&stylist.Services,
		&stylist.Availability,
		&stylist.NoOfCustomerBookings,
		&stylist.NoOfCurrentCustomers,
		&stylist.CreatedAt,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve stylist profile: " + err.Error(),
		})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Stylist profile created successfully",
		"data":    stylist,
	})
}

func ViewStylistProfile(c *fiber.Ctx) error {
	// To get the stylist id from the route parameter
	stylistID := c.Params("stylistId")

	// Struct to hold the combined response
	type StylistProfile struct {
		UserData    models.User    `json:"user"`
		StylistData models.Stylist `json:"stylist"`
	}

	var profile StylistProfile

	// To query user data
	err := config.DB.QueryRow(
		c.Context(),
		"SELECT id, name, email, number, role, location, created_at FROM users WHERE id=$1",
		stylistID,
	).Scan(&profile.UserData.ID, &profile.UserData.Name, &profile.UserData.Email, &profile.UserData.Number, &profile.UserData.Role, &profile.UserData.Location, &profile.UserData.CreatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Stylist not found: " + err.Error(),
		})
	}

	// To query stylist data
	err = config.DB.QueryRow(
		c.Context(),
		"SELECT id, stylist_id, ratings, services, availability, no_of_customer_bookings, no_of_current_customers, created_at FROM stylists WHERE stylist_id=$1",
		stylistID,
	).Scan(&profile.StylistData.ID, &profile.StylistData.StylistID, &profile.StylistData.Ratings, &profile.StylistData.Services, &profile.StylistData.Availability, &profile.StylistData.NoOfCustomerBookings, &profile.StylistData.NoOfCurrentCustomers, &profile.StylistData.CreatedAt)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Stylist data not found: " + err.Error(),
		})
	}

	// To dynamically calculate the number of active customer bookings
	err = config.DB.QueryRow(
		c.Context(),
		`SELECT COUNT(*) FROM bookings 
		 WHERE stylist_id=$1 AND booking_status NOT IN ('completed', 'cancelled')`,
		stylistID,
	).Scan(&profile.StylistData.NoOfCustomerBookings)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to calculate customer bookings: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Stylist profile retrieved successfully",
		"data":    profile,
	})
}

// For the stylist to edit and update their profile
func UpdateStylistProfile(c *fiber.Ctx) error {
	stylistID := c.Params("stylistId")
	// stylistID := c.Locals("user").(float64)

	var input struct {
		Services             *string `json:"services"`
		Availability         *string `json:"availability"`
		NoOfCurrentCustomers *int    `json:"no_of_current_customers"`
	}

	// To parse the req body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input" + err.Error(),
		})
	}

	// To build the update query dynamically based on provided fields
	query := "UPDATE stylists SET "
	args := []interface{}{}
	argCounter := 1

	if input.Services != nil {
		query += "services=$" + strconv.Itoa(argCounter) + ", "
		args = append(args, *input.Services)
		argCounter++
	}

	if input.Availability != nil {
		query += "availability=$" + strconv.Itoa(argCounter) + ", "
		args = append(args, *input.Availability)
		argCounter++
	}

	if input.NoOfCurrentCustomers != nil {
		query += "no_of_current_customers=$" + strconv.Itoa(argCounter) + ", "
		args = append(args, *input.NoOfCurrentCustomers)
		argCounter++
	}

	// Remove trailing comma and space
	query = strings.TrimSuffix(query, ", ")
	query += " WHERE stylist_id=$" + strconv.Itoa(argCounter)
	args = append(args, stylistID)

	// To execute the query
	_, err := config.DB.Exec(c.Context(), query, args...)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update the stylist data" + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"error": "Stylist profile updated successfully",
	})
}

// For the stylist to edit and update their profile
func EditStylistProfile(c *fiber.Ctx) error {
	stylistID := c.Locals("user").(float64)

	// To parse the req body
	var input struct {
		Services     string `json:"services"`
		Availability string `json:"availability"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	// To update the stylist profile
	_, err := config.DB.Exec(
		c.Context(),
		"UPDATE stylists SET services=$1, availability=$2 WHERE stylist_id=$3",
		input.Services, input.Availability, stylistID,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update stylist profile: " + err.Error(),
		})
	}

	// To fetch the updated stylist profile
	var stylist models.Stylist

	err = config.DB.QueryRow(
		c.Context(),
		"SELECT * FROM stylists WHERE stylist_id=$1",
		stylistID,
	).Scan(
		&stylist.ID,
		&stylist.StylistID,
		&stylist.Ratings,
		&stylist.Services,
		&stylist.Availability,
		&stylist.NoOfCustomerBookings,
		&stylist.NoOfCurrentCustomers,
		&stylist.CreatedAt,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve stylist profile: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Stylist profile updated successfully",
		"data":    stylist,
	})
}
