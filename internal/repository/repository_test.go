package repository_test

import (
	"testing"

	"github.com/raihannurr/simple-auth-api/internal/config"
	"github.com/raihannurr/simple-auth-api/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestRepository(t *testing.T) {
	simpleStruct := repository.InitializeRepository(config.AppConfig{})
	mysql := repository.InitializeRepository(config.AppConfig{Database: config.DatabaseConfig{Adapter: "mysql"}})

	assert.IsType(t, &repository.SimpleStructRepository{}, simpleStruct)
	assert.IsType(t, &repository.MysqlRepository{}, mysql)
}
