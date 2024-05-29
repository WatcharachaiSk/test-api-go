package model

import "github.com/watcharachai/fiber-db/internal/database"

func MigrationDB() {
	database.DB.AutoMigrate(&User{})
}
