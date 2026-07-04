package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		Port: getEnv("APP_PORT", "8080"),
		Env:  getEnv("APP_ENV", "dev"),

		PostgresDSN: getEnv("POSTGRES_DSN", ""),
	}

	if cfg.PostgresDSN == "" {
		log.Println("warning: POSTGRES_DSN is empty")
	}

	return cfg
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
