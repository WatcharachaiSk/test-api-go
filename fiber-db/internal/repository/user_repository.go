package repository

import (
	"github.com/watcharachai/fiber-db/internal/database"
	"github.com/watcharachai/fiber-db/internal/model"
)

func GetAllUsers() ([]model.User, error) {
	var users []model.User
	result := database.DB.Find(&users)
	return users, result.Error
}

func CreateUser(user model.User) error {
	result := database.DB.Create(&user)
	return result.Error
}

func FindUserById(id int) (model.User, error) {
	var user model.User
	result := database.DB.Where("id = ?", id).First(&user)
	return user, result.Error
}
