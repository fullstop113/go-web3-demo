package config

import (
	"os"

	"github.com/joho/godotenv"
	"errors"
)


type Config struct {
	HTTPAddr string
	DBDSN string
	JWTSecret string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		HTTPAddr: getEnv("HTTP_ADDR", "8080"),
		DBDSN: getEnv("DB_DSN", "app_db"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
	if len(cfg.JWTSecret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters long")
	}
	return cfg, nil
}


func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}