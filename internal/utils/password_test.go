package utils_test

import (
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	password := "password"
	hashed := utils.HashPassword(password)
	assert.NotEmpty(t, hashed)

	assert.True(t, utils.VerifyPassword(password, hashed))
	assert.False(t, utils.VerifyPassword("wrongpassword", hashed))
}

func TestIsStrongPassword(t *testing.T) {
	assert.True(t, utils.IsStrongPassword("Password123!"))
	assert.False(t, utils.IsStrongPassword("Pwd123!"))
	assert.False(t, utils.IsStrongPassword("password"))
	assert.False(t, utils.IsStrongPassword("Password"))
	assert.False(t, utils.IsStrongPassword("password123"))
}
