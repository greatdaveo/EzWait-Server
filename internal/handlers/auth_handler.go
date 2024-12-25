package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserSerializer struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Number string `json:"number"`
	Role   string `json:"role"`
	// Password        string    `json:"password"`
	// ConfirmPassword string    `json:"confirm_password"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateResponseUser(userModel models.User) UserSerializer {
	return UserSerializer{ID: userModel.ID, Name: userModel.Name, Email: userModel.Email, Number: userModel.Number, Role: userModel.Role, Location: userModel.Location, CreatedAt: userModel.CreatedAt}
}

func RegisterHandler(c *fiber.Ctx) error {
	// To Parse the req body into a User struct
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// To validate the user roles
	if user.Role != models.RoleStylist && user.Role != models.RoleCustomer {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid role. Must be stylist or customer"})
	}

	// To check passwords
	if user.Password != user.ConfirmPassword {
		return c.Status(400).JSON(fiber.Map{"error": "Passwords do not match"})
	}

	// To Hash the Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return c.Status(500).JSON(fiber.Map{"error": "Internal server error"})
	}

	user.Password = string(hashedPassword)

	// To save the user to the DB
	_, err = config.DB.Exec(
		c.Context(),
		"INSERT INTO users (name, email, role, number, password, location, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW())",
		user.Name, strings.ToLower(user.Email), user.Role, user.Number, user.Password, user.Location,
	)

	if err != nil {
		log.Println("Error saving user:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	responseUser := CreateResponseUser(user)

	return c.Status(201).JSON(fiber.Map{"message": "User created successfully", "data": responseUser})

}

func LoginHandler(c *fiber.Ctx) error {
	// To Parse the request body
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid inputs, please enter a valid data"})
	}

	// To fetch user from the DB
	var user models.User
	err := config.DB.QueryRow(
		c.Context(),
		"SELECT id, name, email, password, location FROM users WHERE email=$1",
		input.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Location)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// To compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	responseUser := CreateResponseUser(user)

	// Response
	return c.Status(200).JSON(fiber.Map{
		"message": "Login successful",
		"data":    responseUser,
	})

}
