package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Port        string `port:"PORT"`
	JWTSecret   string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using system env")
	}

	return &Config{
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		Port:        mustGetEnv("PORT"),
		JWTSecret:   mustGetEnv("JWT_SECRET"),
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf(" ENV '%s' need to be filled", key)
	}
	return v
}
