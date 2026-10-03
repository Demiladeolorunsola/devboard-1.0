package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	Database  string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	port := os.Getenv("DEVBOARD_PORT")
	if port == "" {
		port = "8080"
	}

	database := os.Getenv("DEVBOARD_DATABASE")
	if database == "" {
		database = "devboard.db"
	}

	jwtSecret := os.Getenv("DEVBOARD_JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "devboard-development-secret"
	}

	return &Config{
		Port:      port,
		Database:  database,
		JWTSecret: jwtSecret,
	}
}
