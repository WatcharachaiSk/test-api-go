// internal/middleware/auth.go
package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func Auth(c *fiber.Ctx) error {
	// Authentication logic
	return c.Next()
}
