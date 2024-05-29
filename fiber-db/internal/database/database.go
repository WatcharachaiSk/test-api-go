// internal/database/database.go
package database

import (
	"fmt"

	"github.com/watcharachai/fiber-db/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	cfg := config.GetConfig()

	host := cfg.DBHost         // or the Docker service name if running in another container
	port := cfg.DBPort         // default PostgreSQL port
	username := cfg.DBUsername // as defined in docker-compose.yml
	password := cfg.DBPassword // as defined in docker-compose.yml
	dbname := cfg.DBName       // as defined in docker-compose.yml

	// Connection string
	dsn := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: ")
	}

	fmt.Println("Database Successfully connected!")
}
