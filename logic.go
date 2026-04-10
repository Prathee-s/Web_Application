"Handler" 

package main

import (
	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()
	web()
	students()
	authenticate()

	app.Listen(":8080")
}
