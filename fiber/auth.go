package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

// UserData represents the user data extracted from the JWT token
type UserData struct {
	Username string
	Role     string
}

// userContextKey is the key used to store user data in the Fiber context
const userContextKey = "user"

// extractUserFromJWT is a middleware that extracts user data from the JWT token
func extractUserFromJWT(c *fiber.Ctx) error {
	user := &UserData{}

	// Extract the token from the Fiber context (inserted by the JWT middleware)
	token := c.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	fmt.Println(claims)

	user.Username = claims["username"].(string)
	user.Role = claims["role"].(string)

	// Store the user data in the Fiber context
	c.Locals(userContextKey, user)

	return c.Next()
}

func createToken(id int, username string, role string) (string, error) {
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	secretKey := os.Getenv("SECRET")
	if secretKey == "" {
		return "", fmt.Errorf("SECRET environment variable not set")
	}

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id
	claims["username"] = username // correct variable name
	claims["role"] = role         // correct variable name
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
