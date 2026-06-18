package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	MongoDB   string
	JWTSecret string
	Port      string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		MongoDB:   os.Getenv("MONGO_DB"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
	}

	if cfg.MongoURI == "" {
		return nil, fmt.Errorf("MONGO_URI is required")
	}
	if cfg.MongoDB == "" {
		return nil, fmt.Errorf("MONGO_DB is required")
	}
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
