package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func MakeBooking(c *fiber.Ctx) error {
	customerIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	customerID := uint(customerIDFloat)

	// To parse the booking request body
	var input struct {
		StylistID  uint      `json:"stylist_id"`
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
		BookingDay string    `json:"booking_day"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	// To check if the stylist exist
	var stylist models.Stylist
	if err := config.DB.Where("stylist_id = ?", input.StylistID).First(&stylist).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Stylist not found: " + err.Error(),
		})
	}

	// To know if the time slot is available for bookings
	var count int64
	config.DB.Model(&models.Booking{}).Where(
		"stylist_id = ? AND booking_day = ? AND (start_time < ? AND end_time > ?)",
		input.StylistID, input.BookingDay, input.EndTime, input.StartTime,
	).Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Time slot already booked",
		})
	}

	// To handle the booking day format
	bookingDay, err := time.Parse("2006-01-02", input.BookingDay)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid date format" + err.Error(),
		})
	}

	// To handle auto confirm settings
	bookingStatus := "pending"
	if stylist.AutoConfirm {
		bookingStatus = "confirmed"
	}

	// To create the booking with "pending" status
	booking := models.Booking{
		UserID:        customerID,
		StylistID:     input.StylistID,
		StartTime:     input.StartTime,
		EndTime:       input.EndTime,
		BookingDay:    bookingDay,
		BookingStatus: bookingStatus,
		CreatedAt:     time.Now(),
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create booking: " + err.Error(),
		})
	}

	// To increase the no of active bookings for the stylist
	config.DB.Model(&stylist).Update("no_of_customer_bookings", stylist.NoOfCustomerBookings+1)

	return c.Status(201).JSON(fiber.Map{
		"message": "Booking created successfully",
		"data":    booking,
	})
}

func ViewAllBookings(c *fiber.Ctx) error {
	userIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}

	userID := uint(userIDFloat)

	role := c.Locals("role").(string)
	// To query params for filtering
	statusFilter := c.Query("status")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 50 {
		limit = 10
	}

	offset := (page - 1) * limit

	var bookings []models.Booking
	query := config.DB.Model(&models.Booking{})

	if role == "customer" {
		query = query.Where("user_id = ?", userID)
	} else if role == "stylist" {
		query = query.Where("stylist_id = ?", userID)
	}

	if statusFilter != "" {
		query = query.Where("booking_status = ?", statusFilter)
	}

	if err := query.Offset(offset).Limit(limit).Find(&bookings).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch bookings: " + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Bookings received successfully",
		"data":    bookings,
		"page":    page,
		"limit":   limit,
	})
}

func EditBooking(c *fiber.Ctx) error {
	customerIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	customerID := uint(customerIDFloat)

	// To get the booking ID from the req params
	bookingIDStr := c.Params("bookingId")
	// fmt.Println("Raw booking ID string from request:", bookingIDStr)

	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil || bookingID < 1 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid booking ID: " + err.Error(),
		})
	}

	// fmt.Println("Extracted booking ID:", bookingID)

	// To fetch the booking form the DB
	var booking models.Booking
	if err := config.DB.Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Booking not found: " + err.Error(),
		})
	}

	// To check if the bookings belong to the authenticated user
	if booking.UserID != customerID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to edit this booking",
		})
	}

	if booking.BookingStatus != "pending" {
		return c.Status(400).JSON(fiber.Map{
			"error": "You can only edit a pending bookings",
		})
	}

	// To parse the req body
	var input struct {
		StartTime  time.Time `json:"start_time"`
		EndTime    time.Time `json:"end_time"`
		BookingDay string    `json:"booking_day"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	// To convert the booking day to the correct format
	bookingDay, err := time.Parse("2006-01-02", input.BookingDay)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid booking day format" + err.Error(),
		})
	}

	// To check if the new time slot is avaliable
	var count int64
	config.DB.Model(&models.Booking{}).Where(
		"stylist_id = ? AND booking_day = ? AND (start_time < ? AND end_time > ?)",
		booking.StylistID, bookingDay, input.EndTime, input.StartTime,
	).Count(&count)

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Time slot already booked",
		})
	}

	// To update the booking details
	booking.StartTime = input.StartTime
	booking.EndTime = input.EndTime
	booking.BookingDay = bookingDay

	if err := config.DB.Save(&booking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update booking" + err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Booking updated successfully",
		"data":    booking,
	})
}

func UpdateBookingStatus(c *fiber.Ctx) error {
	stylistIDFloat, ok := c.Locals("user").(float64)
	if !ok {
		return c.Status(500).JSON(fiber.Map{
			"error": "Invalid user ID format",
		})
	}
	stylistID := uint(stylistIDFloat)

	// To get booking ID from Params
	bookingIDStr := c.Params("bookingId")

	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil || bookingID < 1 {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid booking ID: " + err.Error(),
		})
	}

	// To get the new status from the req body
	var input struct {
		NewStatus string `json:"new_status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid input: " + err.Error(),
		})
	}

	allowedStatuses := map[string]bool{
		"pending":   true,
		"confirmed": true,
		"completed": true,
		"cancelled": true,
	}

	if !allowedStatuses[input.NewStatus] {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid booking status",
		})
	}

	// To etch the booking from the DB
	var booking models.Booking
	if err := config.DB.Where("id = ?", bookingID).First(&booking).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Booking not found" + err.Error(),
		})
	}

	// To check if the stylist is the owner of this booking
	if booking.StylistID != stylistID {
		return c.Status(403).JSON(fiber.Map{
			"error": "You are not authorized to update this booking",
		})
	}

	// To prevent updates on completed/cancelled bookings
	if booking.BookingStatus == "completed" || booking.BookingStatus == "cancelled" {
		return c.Status(400).JSON(fiber.Map{
			"error": "This booking cannot be updated",
		})
	}

	// To update the booking status
	booking.BookingStatus = input.NewStatus
	if err := config.DB.Save(&booking).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update booking status",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Booking status updated successfully",
		"data":    booking,
	})
}

// To check for expired bookings and automatically mark them as completed
func MarkCompletedBookings() {
	var expiredBookings []models.Booking
	// To find booking where end_time has passed and status is still "confirmed"
	if err := config.DB.Where("end_time < ? AND booking_status = ?", time.Now(), "confirmed").Find(&expiredBookings).Error; err != nil {
		fmt.Println("Error fetching expired bookings: ", err)
		return
	}

	// To update each expired booking
	for _, booking := range expiredBookings {
		booking.BookingStatus = "completed"
		config.DB.Save(&booking)
		fmt.Println("Booking marked as completed: ", booking.ID)
	}
}
