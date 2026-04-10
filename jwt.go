package main

import (
	"time"

	jwtware "github.com/gofiber/contrib/v3/jwt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = "supersecretkey" // keep in env variable in production

func authenticate() {
	app := fiber.New()

	// Public login route
	app.Post("/login", func(c fiber.Ctx) error {
		type LoginRequest struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		req := new(LoginRequest)
		if err := c.Bind().JSON(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		// Demo check
		if req.Username != "admin" || req.Password != "password" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
		}

		// Create JWT token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": req.Username,
			"exp":      time.Now().Add(time.Hour).Unix(),
		})

		t, err := token.SignedString([]byte(secretKey)) // signing needs []byte
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.JSON(fiber.Map{"token": t})
	})

	// Protect /employees routes with JWT middleware
	app.Use("/employees", jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secretKey)}, // middleware expects struct
	}))

	// Protected route
	app.Get("/employees", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "You accessed a protected route!"})
	})
}
