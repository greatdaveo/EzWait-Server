package main

import (
	"ezwait/config"
	"ezwait/internal/routers"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Connect to DB
	config.ConnectDB()

	// config.RunMigrations()
	// config.DB.Exec("ALTER TABLE stylists DROP CONSTRAINT IF EXISTS fk_bookings_stylist;")

	// To call mark completed bookings every minute
	// go func() {
	// 	for {
	// 		handlers.MarkCompletedBookings()
	// 		time.Sleep(1 * time.Hour)
	// 	}
	// }()

	// Fiber app
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// To set up route
	routers.SetupRoutes(app)

	// Start the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
