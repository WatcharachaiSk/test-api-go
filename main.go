package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// DB connection
	dbUsername := os.Getenv("DATABASE_USERNAME")
	dbPassword := os.Getenv("DATABASE_PASSWORD")
	dbName := os.Getenv("DATABASE_DATABASE")
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	runPort := os.Getenv("PORT")

	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", dbHost, dbUsername, dbName, dbPort, dbPassword)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	db.AutoMigrate(User{})

	// Route ----------------------------------------------------------------
	router := gin.Default()

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})

	router.POST("/users", func(c *gin.Context) {
		user := User{}
		if err := c.BindJSON(&user); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		user.Password = string(hashedPassword)
		if err := CreateUser(db, &user); err != nil {
			c.JSON(500, gin.H{"error": "Failed to create user"})
			return
		}
		user.Password = ""
		c.JSON(201, gin.H{"message": "User created successfully", "user": user})
	})
	router.GET("/users", func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch users"})
			return
		}

		// สร้างข้อมูลผู้ใช้ที่ไม่รวมฟิลด์รหัสผ่าน
		var response []gin.H
		for _, user := range users {
			response = append(response, gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			})
		}
		c.JSON(200, response)
	})
	router.GET("/users/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		// สร้างข้อมูลผู้ใช้ที่ไม่รวมฟิลด์รหัสผ่าน
		response := gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}

		c.JSON(200, response)
	})

	// Run server ----------------------------------------------------
	fmt.Println("Server Running on Port: ", runPort)
	router.Run(":" + runPort)
}

func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
