package repository

import (
	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/errors"
	"github.com/raihannurr/simple-auth-api/internal/utils"
)

var _ IRepository = &SimpleStructRepository{} // Interface Check

type SimpleStructRepository struct {
	UserIDIncrement  uint
	UsersByID        map[uint]*entity.User
	UserIDByUsername map[string]uint
	UserIDByEmail    map[string]uint
}

func (r *SimpleStructRepository) CreateUser(username string, email string, password string) (entity.User, error) {
	_, usernameExist := r.UserIDByUsername[username]
	_, emailExist := r.UserIDByEmail[email]
	if usernameExist || emailExist {
		return entity.User{}, errors.ErrUserExists
	}

	r.UserIDIncrement++
	ID := r.UserIDIncrement

	user := entity.User{
		ID:          ID,
		Username:    username,
		Email:       email,
		Password:    utils.HashPassword(password),
		Verified:    false,
		Description: "",
	}

	r.UserIDByUsername[username] = user.ID
	r.UserIDByEmail[email] = user.ID
	r.UsersByID[user.ID] = &user
	return user, nil
}

func (r *SimpleStructRepository) GetUserByID(id uint) (entity.User, error) {
	user, ok := r.UsersByID[id]
	if !ok {
		return entity.User{}, errors.ErrUserNotFound
	}

	return *user, nil
}

func (r *SimpleStructRepository) GetUserByUsername(username string) (entity.User, error) {
	id, ok := r.UserIDByUsername[username]
	if !ok {
		return entity.User{}, errors.ErrUserNotFound
	}

	return *r.UsersByID[id], nil
}

func (r *SimpleStructRepository) UpdateUser(id uint, updatedUser entity.User) (entity.User, error) {
	user, ok := r.UsersByID[id]
	if !ok {
		return entity.User{}, errors.ErrUserNotFound
	}

	user.Description = updatedUser.Description

	return *user, nil
}

func NewSimpleStructRepository() *SimpleStructRepository {
	return &SimpleStructRepository{
		UsersByID:        make(map[uint]*entity.User),
		UserIDByUsername: make(map[string]uint),
		UserIDByEmail:    make(map[string]uint),
		UserIDIncrement:  0,
	}
}
