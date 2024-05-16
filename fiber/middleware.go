package main

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

func checkMiddleware(c *fiber.Ctx) error {
	// Start timer
	start := time.Now()
	// ใช้ฟังก์ชัน Format เพื่อแปลง time.Time เป็นสตริง
	formattedTime := start.Format("2006-01-02 15:04:05")
	fmt.Printf("URL = %s, Method = %s, Time = %s \n", c.OriginalURL(), c.Method(), formattedTime)

	return c.Next()
}
