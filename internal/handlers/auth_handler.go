package handlers

import (
	"ezwait/config"
	"ezwait/internal/models"
	"ezwait/internal/utils"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
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

var store = session.New()

func CreateResponseUser(userModel models.User) UserSerializer {
	return UserSerializer{ID: userModel.ID, Name: userModel.Name, Email: userModel.Email, Number: userModel.Number, Role: userModel.Role, Location: userModel.Location, CreatedAt: userModel.CreatedAt}
}

func RegisterHandler(c *fiber.Ctx) error {
	// To Parse the req body into a User struct
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input - " + err.Error()})
	}

	user.Email = strings.ToLower(user.Email)

	var existingUserID uint
	err := config.DB.QueryRow(
		c.Context(),
		"SELECT id FROM users WHERE email=$1",
		user.Email,
	).Scan(&existingUserID)

	if err == nil {
		return c.Status(400).JSON(fiber.Map{"error": "Email already registered. Please log in."})
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
	user.CreatedAt = time.Now()

	// To save the user to the DB
	err = config.DB.QueryRow(
		c.Context(),
		"INSERT INTO users (name, email, password, number, role, location, created_at) VALUES ($1, $2, $3, $4, $5, $6, NOW()) RETURNING id, created_at",
		user.Name, user.Email, user.Password, user.Number, user.Role, user.Location,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		log.Println("Error saving user:", err)
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	responseUser := CreateResponseUser(user)
	// fmt.Println("Registered user: ", responseUser)

	return c.Status(201).JSON(fiber.Map{"message": "User created successfully", "data": responseUser})

}

func LoginHandler(c *fiber.Ctx) error {
	// To retrieve email and password from the request body
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var loginReq LoginRequest

	// To Parse the incoming request body
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid inputs, please enter a valid data"})
	}

	// To fetch user from the DB
	var user models.User

	err := config.DB.QueryRow(
		c.Context(),
		"SELECT id, name, email, password, number, role, location FROM users WHERE email=$1",
		loginReq.Email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Number, &user.Role, &user.Location)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// To compare the hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid email or password"})
	}

	// To generate JWT
	token, err := utils.GenerateToken(&user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	responseUser := CreateResponseUser(user)
	fmt.Println("Login user: ", responseUser)

	// Response
	return c.Status(200).JSON(fiber.Map{
		"message": "Login successful",
		"token":   token,
		"data":    responseUser,
	})
}

func LogoutHandler(c *fiber.Ctx) error {
	// To retrieve the session
	session, err := store.Get(c)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// To destroy the session
	if err := session.Destroy(); err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Successfully logged out",
	})
}
