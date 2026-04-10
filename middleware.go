package main

import (
	"github.com/gofiber/fiber/v3"
	// "github.com/"
)

type Student struct {
	Name  string `json:"name" validate:"requiered, min=3, max=20"`
	Phone int32  `json:"phone"`
}

func students() {
	app := fiber.New()

	app.Get("/student", func(c fiber.Ctx) error {
		return c.JSON(stu)
	})

	app.Post("/stu_create", func(c fiber.Ctx) error {

		if len(stu) < 3 {
			// if len(stu) < 3 || len(stu) > 25 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name should be Greater than 3 and less than 25..!",
			})
		}

		stu = append(stu, stu...)

		return c.Status(fiber.StatusCreated).JSON(stu)
	})
}
