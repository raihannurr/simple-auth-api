package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port    string
	JWT     JWTConfig
	Session SessionConfig
}

type JWTConfig struct {
	PrivateKey string
	Lifetime   time.Duration
	Issuer     string
}

type SessionConfig struct {
	SecretKey       string
	CSRFTokenLength int
	CSRFTokenExp    time.Duration
}

func LoadEnv(filename string) AppConfig {
	log.Println("Loading .env file...")
	godotenv.Load(filename)

	lifetime, err := strconv.Atoi(os.Getenv("JWT_LIFETIME_IN_MINUTES"))
	if err != nil {
		lifetime = 10 // default lifetime in minutes
	}

	csrfTokenExp, err := strconv.Atoi(os.Getenv("CSRF_TOKEN_EXP_IN_MINUTES"))
	if err != nil {
		csrfTokenExp = 10 // default lifetime in minutes
	}

	csrfTokenLength, err := strconv.Atoi(os.Getenv("CSRF_TOKEN_LENGTH"))
	if err != nil {
		csrfTokenLength = 10 // default length
	}

	return AppConfig{
		Port: os.Getenv("APP_PORT"),
		JWT: JWTConfig{
			PrivateKey: os.Getenv("JWT_PRIVATE_KEY"),
			Lifetime:   time.Duration(lifetime) * time.Minute,
			Issuer:     os.Getenv("JWT_ISSUER"),
		},
		Session: SessionConfig{
			SecretKey:       os.Getenv("SESSION_SECRET_KEY"),
			CSRFTokenLength: csrfTokenLength,
			CSRFTokenExp:    time.Duration(csrfTokenExp) * time.Minute,
		},
	}
}
