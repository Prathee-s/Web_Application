package main

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Employee struct {
	Name  string `json:"name" validate:"required,min=3,max=32"`
	ID    int    `json:"id"`
	Email string `json:"email" validate:"required,email"`
	Age   int    `json:"age" validate:"gte=18,lte=65"`
}

var (
	users    = []Employee{}
	validate = validator.New()
	stu      = []Student{}
	nextID   = 1
)

func validateUser(user *Employee) []string {
	var errs []string
	if err := validate.Struct(user); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, e.Field()+" failed on "+e.Tag())
		}
	}
	return errs
}

func web() {
	app := fiber.New()

	// CREATE (POST)
	app.Post("/users", func(c fiber.Ctx) error {
		user := new(Employee)
		if err := c.Bind().JSON(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
		}

		errs := validateUser(user)
		if len(errs) > 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
		}

		user.ID = nextID
		nextID++
		users = append(users, *user)

		return c.Status(fiber.StatusCreated).JSON(user)
	})

	app.Get("/users", func(c fiber.Ctx) error {
		return c.JSON(users)
	})

	// READ ONE (GET by ID)
	app.Get("/users/:id", func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		for _, u := range users {
			if u.ID == id {
				return c.JSON(u)
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	})

	// UPDATE (PUT)
	app.Put("/users/:id", func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		for i, u := range users {
			if u.ID == id {
				updated := new(Employee)
				if err := c.Bind().JSON(updated); err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid JSON"})
				}

				errs := validateUser(updated)
				if len(errs) > 0 {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"errors": errs})
				}

				updated.ID = id
				users[i] = *updated
				return c.JSON(updated)
			}
		}

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	})

	// DELETE
	app.Delete("/users/:id", func(c fiber.Ctx) error {
		id, _ := strconv.Atoi(c.Params("id"))
		for i, u := range users {
			if u.ID == id {
				users = append(users[:i], users[i+1:]...)
				return c.JSON(fiber.Map{"message": "User deleted"})
			}
		}
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	})
}
