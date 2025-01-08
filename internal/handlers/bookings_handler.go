package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BooingSerializer struct {
	ID            uint      `json:"id"`
	UserID        uint      `json:"user_id"`
	StylistID     uint      `json:"stylist_id"`
	BookingTime   string    `json:"booking_time"`
	BookingDay    string    `json:"booking_day"`
	BookingStatus string    `json:"booking_status"`
	CreatedAt     time.Time `json:"created_at"`
}

func BookingResponse(bookingModel models.Booking) BooingSerializer {
	return BooingSerializer{
		ID:            bookingModel.ID,
		UserID:        bookingModel.UserID,
		StylistID:     bookingModel.StylistID,
		BookingTime:   bookingModel.BookingTime,
		BookingDay:    bookingModel.BookingDay,
		BookingStatus: bookingModel.BookingStatus,
		CreatedAt:     bookingModel.CreatedAt}
}

// For users to make bookings
func MakeBooking(c *fiber.Ctx) error {

	var booking models.Booking

	if err := c.BodyParser(&booking); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if booking.BookingTime == "" || booking.BookingDay == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Booking time and day are required",
		})
	}

	row := config.DB.QueryRow(
		c.Context(),
		"INSERT INTO bookings (user_id, stylist_id, booking_time, booking_day, booking_status, created_at) VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id, created_at",
		booking.UserID, booking.StylistID, booking.BookingTime, booking.BookingDay, booking.BookingStatus,
	)

	if err := row.Scan(&booking.ID, &booking.CreatedAt); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to make booking",
		})
	}

	bookingResponse := BookingResponse(booking)

	fmt.Println("bookingResponse: ", bookingResponse)

	return c.Status(201).JSON(fiber.Map{
		"message": "Booking created successfully",
		"data":    bookingResponse,
	})
}

func EditBooking(c *fiber.Ctx) error {
	id := c.Params("bookingsId")

	var booking models.Booking
	if err := c.BodyParser(&booking); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if booking.BookingTime == "" || booking.BookingDay == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Booking time and day are required",
		})
	}

	// To update the booking in the database
	_, err := config.DB.Exec(
		c.Context(),
		"UPDATE bookings SET booking_time=$1, booking_day=$2 WHERE id=$3",
		booking.BookingTime, booking.BookingDay, id,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update booking",
		})
	}

	// To fetch the updated booking
	var updatedBooking models.Booking
	err = config.DB.QueryRow(
		c.Context(),
		"SELECT id, user_id, stylist_id, booking_time, booking_day, booking_status created_at FROM bookings WHERE id=$1",
		id,
	).Scan(
		&updatedBooking.ID,
		&updatedBooking.UserID,
		&updatedBooking.StylistID,
		&updatedBooking.BookingTime,
		&updatedBooking.BookingDay,
		&updatedBooking.BookingStatus,
		// &updatedBooking.CreatedAt,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to fetch updated booking: " + err.Error(),
		})
	}

	// updatedBooking = BookingResponse(updatedBooking)

	return c.Status(200).JSON(fiber.Map{
		"message": "Booking updated successfully",
		"data":    updatedBooking,
	})

}

func ViewAllBookings(c *fiber.Ctx) error {
	stylistID := c.Params("stylistId")
	rows, err := config.DB.Query(c.Context(), "SELECT * FROM bookings WHERE stylist_id=$1", stylistID)

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
	id := c.Params("bookingsId")
	var input struct {
		Status string `json:"booking_status"`
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

// func UpdateCurrentCustomers(c *fiber.Ctx) error {
// 	var stylist models.Stylist

// 	stylistID := c.Params("id")
// 	var input struct {
// 		Action string `json:"action"`
// 	}

// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	query := ""
// 	if input.Action == "increment" {
// 		query = "UPDATE stylists SET no_of_current_customers = no_of_current_customers + 1 WHERE id=$1"
// 	} else if input.Action == "decrement" {
// 		query = "UPDATE stylists SET no_of_current_customers = GREATEST(no_of_current_customers - 1, 0) WHERE id=$1"
// 	} else {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error": "Invalid action",
// 		})
// 	}

// 	_, err := config.DB.Exec(c.Context(), query, stylistID)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{
// 			"error": err.Error(),
// 		})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"message": "Current customers updated successfully",
// 	})
// }
