package repository

import (
	"log"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
)

type IRepository interface {
	CreateUser(username string, email string, password string) (entity.User, error)
	GetUserByID(id uint) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	UpdateUser(id uint, updatedUser entity.User) (entity.User, error)
}

func InitializeRepository(cfg config.AppConfig) IRepository {
	log.Println("Initializing repository with adapter:", cfg.Database.Adapter)
	if cfg.Database.Adapter == "mysql" {
		return NewMysqlRepository(cfg.Database)
	}
	return NewSimpleStructRepository()
}
