package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL  string `env:"DATABASE_URL"`
	Port         string `port:"PORT"`
	JWTSecret    string
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
	AppURL       string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env not found, using system env")
	}

	return &Config{
		DatabaseURL:  mustGetEnv("DATABASE_URL"),
		Port:         mustGetEnv("PORT"),
		JWTSecret:    mustGetEnv("JWT_SECRET"),
		SMTPHost:     mustGetEnv("SMTP_HOST"),
		SMTPPort:     mustGetEnv("SMTP_PORT"),
		SMTPEmail:    mustGetEnv("SMTP_EMAIL"),
		SMTPPassword: mustGetEnv("SMTP_PASSWORD"),
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf(" ENV '%s' need to be filled", key)
	}
	return v
}
