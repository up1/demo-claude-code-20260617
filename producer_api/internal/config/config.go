package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	RedisAddr    string
	RedisPassword string
	RabbitMQURI  string
	DedupTTL     time.Duration
	ServiceName  string
	OTLPEndpoint string // optional; when empty, traces are exported to stdout
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		Port:          os.Getenv("PORT"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RabbitMQURI:   os.Getenv("RABBITMQ_URI"),
		ServiceName:   os.Getenv("SERVICE_NAME"),
		OTLPEndpoint:  os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
	}

	if cfg.RedisAddr == "" {
		return nil, fmt.Errorf("REDIS_ADDR is required")
	}
	if cfg.RabbitMQURI == "" {
		return nil, fmt.Errorf("RABBITMQ_URI is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "producer-api"
	}

	cfg.DedupTTL = 24 * time.Hour
	if raw := os.Getenv("DEDUP_TTL"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("DEDUP_TTL must be a valid duration (e.g. 24h): %w", err)
		}
		cfg.DedupTTL = d
	}

	return cfg, nil
}
