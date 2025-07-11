package config

import "os"

type Config struct {
	Port string
	DB   string
}

func Load() Config {
	return Config{
		Port: getEnv("PORT", "8080"),
		DB:   getEnv("DATABASE_URL", "postgres://postgres:super@localhost:5432/booru?sslmode=disable"),
	}
}

func getEnv(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
