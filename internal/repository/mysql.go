package repository

import (
	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/errors"
	"github.com/raihannurr/simple-auth-api/internal/utils"

	"gorm.io/gorm"
)

var _ IRepository = &MysqlRepository{} // Interface Check

type MysqlRepository struct {
	DB *gorm.DB
}

func (r *MysqlRepository) CreateUser(username string, email string, password string) (entity.User, error) {
	user := entity.User{
		Username:    username,
		Email:       email,
		Password:    utils.HashPassword(password),
		Verified:    false,
		Description: "",
	}

	execution := r.DB.Create(&user)
	if execution.Error != nil {
		return entity.User{}, execution.Error
	}

	return user, nil
}

func (r *MysqlRepository) GetUserByID(id uint) (entity.User, error) {
	var user entity.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return entity.User{}, errors.ErrUserNotFound
	}

	return user, nil
}

func (r *MysqlRepository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return entity.User{}, errors.ErrUserNotFound
	}

	return user, nil
}

func (r *MysqlRepository) UpdateUser(id uint, updatedUser entity.User) (entity.User, error) {
	updateadbleFields := []string{"description"}
	result := r.DB.Model(&updatedUser).Select(updateadbleFields).Updates(updatedUser)
	if result.Error != nil {
		return entity.User{}, result.Error
	}

	return updatedUser, nil
}

func NewMysqlRepository(cfg config.DatabaseConfig) *MysqlRepository {
	db, err := gorm.Open(cfg.Connection, &cfg.GormConfig)
	utils.PanicIfError(err)

	return &MysqlRepository{
		DB: db,
	}
}
