package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	BaseURL     string
	Port        string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	baseURL := os.Getenv("BASE_URL")
	port := os.Getenv("PORT")

	if dbURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}

	if baseURL == "" {
		return nil, errors.New("BASE_URL is required")
	}

	if port == "" {
		return nil, errors.New("PORT is required")
	}

	return &Config{
		DatabaseURL: dbURL,
		BaseURL:     baseURL,
		Port:        port,
	}, nil
}
