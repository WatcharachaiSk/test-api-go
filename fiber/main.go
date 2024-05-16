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
	db_prot := os.Getenv("DATABASE_PORT")
	// JWT Secret Key
	println(db_prot)

	// Fiber
	app := fiber.New()
	// Apply CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Adjust this to be more restrictive if needed
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post("/login", login)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("SECRET")),
	}))

	userGroups := app.Group("/user")

	// Middleware to extract user data from JWT
	userGroups.Use(extractUserFromJWT)
	userGroups.Get("/", getUsers)
	userGroups.Get("/:id", getUser)
	// Setup routes
	userGroups.Use(isAdmin)
	userGroups.Post("/", createUser)
	userGroups.Put("/:id", updateUser)
	userGroups.Delete("/:id", deleteUser)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Listen(":8080")
}
