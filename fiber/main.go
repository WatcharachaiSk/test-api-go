package main

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func main() {
	app := fiber.New()

	users = append(users, User{ID: 1, Username: "rootUser1", Password: "rootPassword1"})
	users = append(users, User{ID: 2, Username: "rootUser2", Password: "rootPassword2"})

	app.Get("/users", getUsers)
	app.Get("/user/:id", getUser)
	app.Post("/user", createUser)
	app.Put("/user/:id", updateUser)
	app.Delete("/user/:id", deleteUser)

	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello World")
	})

	app.Listen(":8080")
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
