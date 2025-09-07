package main

import (
	"ezwait/config"
	"ezwait/internal/routers"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	// Connect to DB
	config.ConnectDB()

	// To initialize Stripe
	config.InitStripe()

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
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	// For security middleware
	app.Use(helmet.New())
	app.Use(compress.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// For rate limiting - 100 requests per minute per IP
	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
	}))

	// For Request logging
	// app.Use(logger.New(logger.Config{
	// 	Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
	// 	TimeFormat: "02-Jan-2006 15:04:05",
	// 	Output:     log.Writer(),
	// }))

	// // For panic recovery
	// app.Use(recover.New())

	// To set up route
	routers.SetupRoutes(app)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(app.Listen(":" + port))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
		"code":    code,
	})
}
