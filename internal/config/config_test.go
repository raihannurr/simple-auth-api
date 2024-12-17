package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/raihannurr/simple-auth-api/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadEnv(t *testing.T) {
	file, err := os.Create(".env.dummy")
	assert.Nil(t, err)
	file.WriteString("JWT_PRIVATE_KEY=secret-token\n")
	file.WriteString("JWT_LIFETIME_IN_MINUTES=\n")
	file.WriteString("JWT_ISSUER=simple-auth-api\n")
	file.WriteString("SESSION_SECRET_KEY=secret-session\n")
	file.WriteString("APP_PORT=9999\n")
	file.Close()
	defer os.Remove(".env.dummy")

	cfg := config.LoadEnv(".env.dummy")
	assert.Equal(t, cfg.JWT.PrivateKey, "secret-token")
	assert.Equal(t, cfg.JWT.Lifetime, 10*time.Minute)
	assert.Equal(t, cfg.JWT.Issuer, "simple-auth-api")
	assert.Equal(t, cfg.Session.SecretKey, "secret-session")
	assert.Equal(t, cfg.Port, "9999")
}
