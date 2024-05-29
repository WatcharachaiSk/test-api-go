package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/watcharachai/fiber-db/internal/handler"
)

func SetupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(("Hello, world!"))
	})
	// Setup routes

	// User
	app.Get("/users", handler.GetUsers)
	app.Get("/user/:id", handler.FindUserById)
	app.Post("/user", handler.CreateUser)
	app.Patch("/user/:id", handler.CreateUser)
}
