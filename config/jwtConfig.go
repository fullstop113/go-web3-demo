package config

import (
	"log"
	"os"
	"github.com/joho/godotenv"
)

func LoadJWTSecret() string {
	_ = godotenv.Load()

	secret := os.Getenv("JWT_SECRET")
	if len(secret) < 32 {
		log.Fatal("JWT_SECRET is missing or too short (min 32 chars)")
	}
	return secret
}
