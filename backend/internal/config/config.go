package config

import "os"

type Config struct {
	PostgresDSN string
}

func MustLoad() Config {
	return Config{
		PostgresDSN: getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/life?sslmode=disable"),
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
