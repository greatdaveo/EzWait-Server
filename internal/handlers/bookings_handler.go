package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"

	"github.com/gofiber/fiber/v2"
)

type BooingSerializer struct {
	ID            uint   `json:"id"`
	UserID        uint   `json:"user_id"`
	StylistID     uint   `json:"stylist_id"`
	BookingTime   string `json:"booking_time"`
	BookingDay    string `json:"booking_day"`
	BookingStatus string `json:"booking_status"`
	CreatedAt     string `json:"created_at"`
}

func BookingResponse(bookingModel models.Booking) BooingSerializer {
	return BooingSerializer{ID: bookingModel.ID, StylistID: bookingModel.StylistID, BookingTime: bookingModel.BookingTime, BookingStatus: bookingModel.BookingStatus, CreatedAt: bookingModel.CreatedAt}
}

func MakeBooking(c *fiber.Ctx) error {
	var booking models.Booking

	if err := c.BodyParser(&booking); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, err := config.DB.Exec(
		c.Context(),
		"INSERT INTO bookings (user_id, stylist_id, booking_time, booking_day, booking_status) VALUES ($1, $2, $3, $4, $5)",
		booking.UserID, booking.StylistID, booking.BookingTime, booking.BookingDay, "pending",
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	bookingResponse := BookingResponse(booking)

	return c.Status(201).JSON(fiber.Map{
		"message": "Booking created successfully",
		"data":    bookingResponse,
	})
}

func EditBooking(c *fiber.Ctx) error {
	id := c.Params("id")

	var booking models.Booking
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, err := config.DB.Exec(
		c.Context(),
		"UPDATE bookings SET booking_time=$1, booking_day=$2 WHERE id=$3",
		booking.BookingTime, booking.BookingDay, id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	bookingResponse := BookingResponse(booking)

	return c.Status(200).JSON(fiber.Map{
		"message": "Booking updated successfully",
		"data":    bookingResponse,
	})

}

func ViewALlBookings(c *fiber.Ctx) error {
	stylistID := c.Params("stylistId")
	rows, err := config.DB.Query(c.Context(), "SELECT * FROM booking WHERE stylist_id=$1", stylistID)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	defer rows.Close()
	var bookings []models.Booking
	for rows.Next() {
		var booking models.Booking
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.StylistID, &booking.BookingTime, &booking.BookingDay, &booking.BookingStatus, &booking.CreatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		bookings = append(bookings, booking)
	}

	return c.Status(200).JSON(bookings)
}

func UpdateBookingStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var input struct {
		Status string `json:"status"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	_, err := config.DB.Exec(
		c.Context(),
		"UPDATE bookings SET booking_status=$1 WHERE id=$2",
		input.Status, id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Booking status updated successfully",
	})
}

func UpdateCurrentCustomers(c *fiber.Ctx) error {
	stylistID := c.Params("id")
	var input struct {
		Action string `json:"action"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	query := ""
	if input.Action == "increment" {
		query = "UPDATE stylists SET no_of_current_customers = no_of_current_customers + 1 WHERE id=$1"
	} else if input.Action == "decrement" {
		query = "UPDATE stylists SET no_of_current_customers = no_of_current_customers - 1 WHERE id=$1"
	} else {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid action",
		})
	}

	_, err := config.DB.Exec(c.Context(), query, stylistID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Current customers updated successfully",
	})
}
