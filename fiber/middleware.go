package main

import (
	"github.com/gofiber/fiber/v2"
)

// isAdmin checks if the user is an admin
func isAdmin(c *fiber.Ctx) error {
	user := c.Locals(userContextKey).(*UserData)

	if user.Role != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}
