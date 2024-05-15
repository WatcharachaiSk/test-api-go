package main

import (
	"fmt"
	"log"
	"os"
	"test-api-go/service"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
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
type RequestLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	// test ----------------------------------------------------------------
	serviceTest := service.TestService()
	fmt.Println(serviceTest)
	// --------------------------------------------------------------------

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
	// Middleware สำหรับตรวจสอบ token

	router.GET("/hello", func(c *gin.Context) {
		fmt.Println("Hello")
		c.JSON(200, gin.H{
			"message": "Hello, Gin!",
		})
	})
	// เส้นทาง GET ที่ต้องตรวจสอบ token

	router.POST("/user", func(c *gin.Context) {
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
		response := gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		}
		c.JSON(201, gin.H{"message": "User created successfully", "user": response})
	})

	router.GET("/users", func(c *gin.Context) {
		// อ่าน claims จาก Context
		c.Get("claims")

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

	router.GET("/user/:id", func(c *gin.Context) {
		var user User
		id := c.Param("id")
		if err := db.First(&user, id).Error; err != nil {

			c.JSON(404, gin.H{"error": fmt.Sprintf("User id is %s not found", id)})
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

	router.PATCH("/user/:id", func(c *gin.Context) {
		user := User{}
		id := c.Param("id")

		// ค้นหาผู้ใช้จาก ID
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		// รับข้อมูลการอัปเดตจากผู้ใช้
		updateData := User{}
		if err := c.BindJSON(&updateData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		if updateData.Password != "" {
			// Hash password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(500, gin.H{"error": "Failed to hash password"})
				return
			}
			updateData.Password = string(hashedPassword)
		}

		// อัปเดตข้อมูลผู้ใช้
		db.Model(&user).Updates(updateData)

		c.JSON(200, gin.H{"message": "User updated successfully"})
	})
	router.DELETE("/user/:id", func(c *gin.Context) {
		user := User{}
		id := c.Param("id")

		// ค้นหาผู้ใช้จาก ID
		if err := db.First(&user, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}
		// Delete user by ID
		if err := db.Delete(&User{}, id).Error; err != nil {
			// Handle error
			c.JSON(500, gin.H{"error": "Failed to delete user"})
			return
		}

		c.JSON(200, gin.H{"message": "User deleted successfully"})
	})

	router.POST("/auth/login", func(c *gin.Context) {
		user := User{}

		body := RequestLogin{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}

		// ค้นหาผู้ใช้จาก Username
		if err := db.Where("username = ?", body.Username).First(&user).Error; err != nil {
			// Handle error
			c.JSON(404, gin.H{"error": "User not found"})
			return
		}

		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
		if err != nil {
			c.JSON(400, gin.H{"error": "Failed password"})
			return
		}

		token, err := createToken(int(user.ID))
		if err != nil {
			// ไม่สามารถสร้าง Token ได้
			panic(err)
		}

		response := gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"token":    token,
		}

		// Found user
		c.JSON(200, response)
	})

	// Run server ----------------------------------------------------
	fmt.Println("Server Running on Port: ", runPort)
	router.Run(":" + runPort)
}

// Func User -------------------------------------------------------------
func CreateUser(db *gorm.DB, user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// JWT ----------------------------------------------------------------
var secretKey = []byte("secret-Tset-Go") // กำหนดคีย์ลับสำหรับเซ็น JWT

func createToken(userID int) (string, error) {
	// สร้างข้อมูลสำหรับ JWT
	token := jwt.New(jwt.SigningMethodHS256)

	// กำหนดเอกสาร (claim) ใน JWT
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // กำหนดเวลาหมดอายุของ Token เป็น 24 ชั่วโมง

	// ลายเซ็น JWT ด้วยคีย์ลับ
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// func parseToken(tokenString string) (jwt.MapClaims, error) {
// 	// ตรวจสอบ Token และแปลงเป็น claim
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		return secretKey, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	// ตรวจสอบว่า Token ถูกต้องหรือไม่
// 	if !token.Valid {
// 		return nil, jwt.ErrSignatureInvalid
// 	}

// 	// แปลง claim ให้เป็น MapClaims
// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return nil, jwt.ErrInvalidKeyType
// 	}

// 	return claims, nil
// }
