package config

import (
	"ezwait/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		"require", // 🔹 Change "require" to "disable" if testing locally
	)

	// var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to the database successfully")
}

func RunMigrations() {
	err := DB.AutoMigrate(&models.User{}, &models.Stylist{}, &models.Booking{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Migrations completed successfully")
}
