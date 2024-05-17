package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/watcharachai/fiber-db/internal/handler"
)

func SetupRoutes(app *fiber.App) {
	// Setup routes
	app.Get("/users", handler.GetUsers)
}
