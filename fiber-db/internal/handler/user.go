// internal/handler/user.go
package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/watcharachai/fiber-db/internal/model"
	"github.com/watcharachai/fiber-db/internal/repository"
	"github.com/watcharachai/fiber-db/internal/utils"
)

func GetUsers(c *fiber.Ctx) error {
	users, err := repository.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to fetch users",
		})
	}
	return c.JSON(users)
}

func CreateUser(c *fiber.Ctx) error {
	user := new(model.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to hash password",
		})
	}
	user.Password = hashedPassword

	if err := repository.CreateUser(*user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to create user",
		})
	}
	return c.JSON(user)
}

func FindUserById(c *fiber.Ctx) error {

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, err := repository.FindUserById(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}
