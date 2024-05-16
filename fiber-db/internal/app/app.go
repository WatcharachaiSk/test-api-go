// internal/app/app.go
package app

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/watcharachai/fiber-db/internal/database"
	"github.com/watcharachai/fiber-db/pkg/config"
)

func Start() {
	cfg := config.GetConfig()
	database.Connect()

	app := fiber.New()
	Port := fmt.Sprintf(":%d", cfg.Port)
	log.Fatal(app.Listen(Port))
}