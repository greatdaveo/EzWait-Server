package routers

import (
	"ezwait/internal/handlers"
	"ezwait/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	// To Register
	api.Post("/register", middleware.ValidateUser, handlers.RegisterHandler)

	// Login User
	api.Post("/login", handlers.LoginHandler)
}
