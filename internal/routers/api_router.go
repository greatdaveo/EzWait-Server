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

	// For Bookings
	api.Post("/customer/bookings", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.MakeBooking)
	api.Put("/bookings/:id", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.EditBooking)
	api.Get("/stylists/:stylistId/bookings", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.ViewALlBookings)
	app.Patch("/bookings/:id/status", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.UpdateBookingStatus)
	app.Patch("/stylists/:id/customers", middleware.AuthMiddleware, middleware.ValidateCustomer, handlers.UpdateCurrentCustomers)
}
