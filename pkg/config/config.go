package config

import (
	"os"
	"time"
)

type Config struct {
	Host string
	Port string
	DB   string

	ReadHeaderTimeout time.Duration
	GracefulTimeout   time.Duration
}

func Load() Config {
	return Config{
		Host: getEnv("HOST", "localhost"),
		Port: getEnv("PORT", "8080"),
		DB:   getEnv("DATABASE_URL", "postgres://postgres:super@localhost:5432/booru?sslmode=disable"),

		ReadHeaderTimeout: 60,
		GracefulTimeout:   8,
	}
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
