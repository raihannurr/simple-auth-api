package entity_test

import (
	"encoding/json"
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	user := entity.User{
		Username: "test",
		Password: "Password123!",
		Email:    "test@example.com",
	}

	json, err := json.Marshal(user)
	assert.NoError(t, err)
	assert.JSONEq(t, `{"id":0,"username":"test","email":"test@example.com","verified":false,"description":""}`, string(json))
}
