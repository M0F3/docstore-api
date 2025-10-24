package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	AppDatabaseConnectionUrl string
	AppAdminDatabaseConnectionUrl string
	JWTSecret                string
}

func LoadAppConfig() *AppConfig {
	_ = godotenv.Load()

	dbConnectionUrl := os.Getenv("APP_DATABASE_CONNECTION_URL")
	if dbConnectionUrl == "" {
		log.Fatal("APP_DATABASE_CONNECTION_URL not set")
	}

	dbAdminConnectionUrl := os.Getenv("APP_ADMIN_DATABASE_CONNECTION_URL")
	if dbConnectionUrl == "" {
		log.Fatal("APP_ADMIN_DATABASE_CONNECTION_URL not set")
	}

	jwtSecret := os.Getenv("APP_JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("APP_JWT_SECRET not set")
	}

	return &AppConfig{
		AppDatabaseConnectionUrl: dbConnectionUrl,
		JWTSecret:                jwtSecret,
		AppAdminDatabaseConnectionUrl: dbAdminConnectionUrl,
	}
}
