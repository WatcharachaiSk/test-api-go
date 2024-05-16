// cmd/myapp/main.go
package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/watcharachai/fiber-db/internal/app"
)

func main() {
	// ENV variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app.Start()
}
