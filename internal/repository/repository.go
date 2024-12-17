package repository

import (
	"github.com/raihannurr/simple-auth-api/internal/entity"
)

type IRepository interface {
	CreateUser(username string, email string, password string) (entity.User, error)
	GetUserByID(id uint) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	UpdateUser(id uint, updatedUser entity.User) (entity.User, error)
}
