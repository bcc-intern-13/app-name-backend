package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string `env:"DATABASE_URL"`
	Port        string `port:"PORT"`

	//supabase object storage
	SupabaseServiceRoleKey string
	StorageBucketCV        string
	StorageBucketAvatar    string
	SupabaseURL            string

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
		DatabaseURL: mustGetEnv("DATABASE_URL"),
		Port:        mustGetEnv("PORT"),

		//supabase object storage
		SupabaseURL:            mustGetEnv("SUPABASE_URL"),
		SupabaseServiceRoleKey: mustGetEnv("SUPABASE_SERVICE_ROLE_KEY"),
		StorageBucketCV:        mustGetEnv("STORAGE_BUCKET_CV"),
		StorageBucketAvatar:    mustGetEnv("STORAGE_BUCKET_AVATAR"),

		//jwt
		JWTSecret: mustGetEnv("JWT_SECRET"),

		//email host
		SMTPHost:     mustGetEnv("SMTP_HOST"),
		SMTPPort:     mustGetEnv("SMTP_PORT"),
		SMTPEmail:    mustGetEnv("SMTP_EMAIL"),
		SMTPPassword: mustGetEnv("SMTP_PASSWORD"),
		AppURL:       mustGetEnv("APP_URL"),
	}
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf(" ENV '%s' need to be filled", key)
	}
	return v
}
