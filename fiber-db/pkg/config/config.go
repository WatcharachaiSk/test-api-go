// pkg/config/config.go
package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	DBUsername string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     int
	Port       int
}

func GetConfig() *Config {
	DBPort, err := strconv.Atoi(os.Getenv("DATABASE_PORT"))
	if err != nil {
		log.Fatal("DBPort ", err)
	}
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("prot ", err)
	}

	return &Config{
		DBUsername: os.Getenv("DATABASE_USERNAME"),
		DBPassword: os.Getenv("DATABASE_PASSWORD"),
		DBName:     os.Getenv("DATABASE_DATABASE"),
		DBHost:     os.Getenv("DATABASE_HOST"),
		DBPort:     DBPort,
		Port:       port,
	}
}
