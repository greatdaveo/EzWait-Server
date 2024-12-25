package main

import (
	"ezwait/config"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// To connect to the database
	config.ConnectDB()
	defer config.DB.Close()

	// To initialize the Fiber app
	app := fiber.New()
	// To set up routes
	// routers.SetupRoutes(app)
	// To start the server
	port := os.Getenv("APP_PORT")
	log.Fatal(app.Listen(":" + port))
}
