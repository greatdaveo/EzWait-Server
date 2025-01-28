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
	// To connect to the database
	config.ConnectDB()
	defer config.DB.Close()

	// To initialize the Fiber app
	app := fiber.New()

	// To add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	// To pass the session store to middleware
	// To set up routes
	routers.SetupRoutes(app)

	// To start the server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
