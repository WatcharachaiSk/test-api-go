package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ReqUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User = []User{
	{ID: 1, Username: "rootUser1", Password: "rootPassword1"},
	{ID: 2, Username: "rootUser2", Password: "rootPassword2"},
}

func login(c *fiber.Ctx) error {
	reqUser := new(ReqUser)

	if err := c.BodyParser(reqUser); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for _, user := range users {
		if user.Username == reqUser.Username && user.Password == reqUser.Password {
			token, err := createToken(user.ID, user.Username)
			println(token)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			return c.JSON(fiber.Map{"token": token})
		}
		if user.Username != reqUser.Username || user.Password != reqUser.Password {
			return fiber.ErrUnauthorized
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("user not found")
}

func createToken(id int, username string) (string, error) {
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
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token
	t, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return t, nil
}
func getUsers(c *fiber.Ctx) error {
	return c.JSON(users)
}
func getUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	for _, user := range users {
		if user.ID == userId {
			return c.JSON(user)
		}
	}
	return c.Status(fiber.StatusNotFound).SendString("user not found")
}
func createUser(c *fiber.Ctx) error {
	user := new(User)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	users = append(users, *user)

	return c.JSON(users)

}
func updateUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	userUpdate := new(User)
	if err := c.BodyParser(userUpdate); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, user := range users {
		if user.ID == userId {
			users[i].Username = userUpdate.Username
			users[i].Password = userUpdate.Password
			return c.JSON(users[i])
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("user not found")
}
func deleteUser(c *fiber.Ctx) error {
	userId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	for i, user := range users {
		if user.ID == userId {
			users = append(users[:i], users[i+1:]...)
			return c.SendStatus(fiber.StatusNoContent)
		}
	}

	return c.Status(fiber.StatusNotFound).SendString("user not found")
}
