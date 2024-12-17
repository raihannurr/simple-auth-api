package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppConfig struct {
	Port     string
	JWT      JWTConfig
	Session  SessionConfig
	Database DatabaseConfig
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

type DatabaseConfig struct {
	Adapter    string
	Connection gorm.Dialector
	GormConfig gorm.Config
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

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

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
		Database: DatabaseConfig{
			Adapter:    os.Getenv("DB_ADAPTER"),
			Connection: mysql.Open(dsn),
			GormConfig: gorm.Config{
				Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
					SlowThreshold:             0,
					LogLevel:                  logger.Info,
					IgnoreRecordNotFoundError: false,
				}),
			},
		},
	}
}
