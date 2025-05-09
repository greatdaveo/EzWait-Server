package routers

import (
	"ezwait/internal/handlers"
	"ezwait/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to EzWait App")
	})

	// For Authentication
	api.Post("/user/register", middleware.ValidateUser, handlers.RegisterHandler)
	api.Post("/user/login", handlers.LoginHandler)
	api.Post("/user/logout", handlers.LogoutHandler)

	// To test session
	app.Get("/test-session", func(c *fiber.Ctx) error {

		user := c.Locals("user")
		role := c.Locals("role")

		return c.JSON(fiber.Map{
			"user": user,
			"role": role,
		})
	})

	// For User Bookings
	api.Post("/customer/bookings", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.MakeBooking)
	api.Get("/view-all/bookings", middleware.AuthMiddleware, handlers.ViewAllBookings)
	api.Get("/view/bookings/:bookingId", middleware.AuthMiddleware, handlers.ViewSingleBooking)
	// For user to edit and update details
	api.Put("/user/edit", middleware.AuthMiddleware, handlers.UpdateUserProfile)

	api.Put("/customer/edit/bookings/:bookingId", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.EditBooking)
	api.Get("/customer/view/all-stylists/", middleware.AuthMiddleware, handlers.ViewAllStylists)

	// api.Patch("/stylists/:id/customers", middleware.AuthMiddleware, handlers.UpdateCurrentCustomers)

	// Stylist Bookings Profile
	// api.Get("/stylists/:stylistId/bookings", middleware.AuthMiddleware, middleware.ValidateStylist, handlers.ViewAllBookings)

	// Stylist
	api.Post("/stylists/profile", middleware.AuthMiddleware, middleware.ValidateStylist, handlers.CreateStylistProfile)
	api.Get("/stylists/:stylistId/profile", middleware.AuthMiddleware, handlers.ViewStylistProfile)
	api.Patch("/bookings/:bookingId/status", middleware.AuthMiddleware, middleware.ValidateStylist, handlers.UpdateBookingStatus)
	api.Patch("/stylists/:stylistId", middleware.AuthMiddleware, middleware.ValidateStylist, handlers.UpdateStylistProfile)
}
