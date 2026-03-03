package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using system env")
	}
	return &Config{
		DatabaseURL: mustGetEnv("DATABASE_URL"),
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf(" ENV '%s' need to be filled", key)
	}
	return v
}
