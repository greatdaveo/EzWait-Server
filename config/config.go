package config

import (
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
	// To only load .env file in local development
	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("SSL_MODE"),
	)

	// var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: false,
		Logger:      logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Connected to the database successfully")
}

func RunMigrations() {
	// DB.Migrator().DropTable(
	// 	&models.User{},
	// 	&models.Stylist{},
	// 	&models.Booking{},
	// )

	// err := DB.AutoMigrate(
	// 	&models.User{},
	// 	&models.Stylist{},
	// 	&models.Booking{},
	// )

	// if err != nil {
	// 	log.Fatal("Failed to migrate database:", err)
	// }

	fmt.Println("Migrations completed successfully")
}
