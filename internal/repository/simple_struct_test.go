package repository_test

import (
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/raihannurr/simple-auth-api/internal/errors"
	"github.com/raihannurr/simple-auth-api/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestSimpleStructRepository(t *testing.T) {
	repo := repository.NewSimpleStructRepository()

	user, err := repo.CreateUser("user1", "user1@example.com", "password")
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, uint(1), user.ID)

	user, err = repo.GetUserByUsername("user1")
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, uint(1), user.ID)

	user, err = repo.GetUserByID(1)
	assert.Nil(t, err)
	assert.Equal(t, "user1", user.Username)
	assert.Equal(t, uint(1), user.ID)

	assert.Equal(t, "", user.Description)
	user, err = repo.UpdateUser(1, entity.User{Description: "user1 description"})
	assert.Nil(t, err)
	assert.Equal(t, "user1 description", user.Description)

	_, err = repo.GetUserByID(2)
	assert.Equal(t, errors.ErrUserNotFound, err)

	_, err = repo.GetUserByUsername("user2")
	assert.Equal(t, errors.ErrUserNotFound, err)

	_, err = repo.UpdateUser(2, entity.User{Description: "user2 description"})
	assert.Equal(t, errors.ErrUserNotFound, err)

	_, err = repo.CreateUser("user1", "", "")
	assert.NotNil(t, err)

	userNew, err := repo.CreateUser("user-new", "user-new@example.com", "password")
	assert.Nil(t, err)
	assert.Equal(t, "user-new", userNew.Username)
	assert.Equal(t, uint(2), userNew.ID)

}
