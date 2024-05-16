package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/joho/godotenv"
)

func main() {
	// ENV variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	prot := os.Getenv("DATABASE_PORT")
	// JWT Secret Key
	println(prot)

	// Fiber
	app := fiber.New()
	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/login", login)
	// Setup routes
	app.Use(checkMiddleware)
	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
	}))
	app.Get("/users", getUsers)
	app.Get("/user/:id", getUser)
	app.Post("/user", createUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Listen(":8080")
}
