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

	// To Register
	api.Post("/user/register", middleware.ValidateUser, handlers.RegisterHandler)
	// Login User
	api.Post("/user/login", handlers.LoginHandler)
	// Logout user
	api.Post("/user/logout", handlers.LogoutHandler)
}
