package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	//databae url
	DatabaseURL string `env:"DATABASE_URL"`
	//app port
	Port string `port:"PORT"`

	//supabase object storage
	SupabaseServiceRoleKey string
	StorageBucketCV        string
	StorageBucketAvatar    string
	SupabaseURL            string

	//jwt secret string
	JWTSecret string

	//gmail configs
	SMTPHost     string
	SMTPPort     string
	SMTPEmail    string
	SMTPPassword string
	AppURL       string
	AppURLFe     string

	//gemini api key
	GeminiAPIKey string

	// midtrans api key and is in production or not
	XenditSecretKey    string
	XenditWebhookToken string

	//google oauth
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string

	//redis
	REDISPassword string
	REDISHost     string
	REDISPort     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("env not found, using system env")
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
		AppURLFe:     mustGetEnv("APP_URL_FE"),

		//gemini api key
		GeminiAPIKey: mustGetEnv("GEMINI_API_KEY"),

		// xendit api key and webhook token
		XenditSecretKey:    mustGetEnv("XENDIT_SECRET_KEY"),
		XenditWebhookToken: mustGetEnv("XENDIT_WEBHOOK_TOKEN"),

		//google oauth
		GoogleClientID:     mustGetEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: mustGetEnv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  mustGetEnv("GOOGLE_REDIRECT_URL"),

		//redis
		REDISPassword: mustGetEnv("REDIS_PASSWORD"),
		REDISHost:     mustGetEnv("REDIS_HOST"),
		REDISPort:     mustGetEnv("REDIS_PORT"),
	}
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal("env value is missing")
		log.Fatalf(" ENV '%s' need to be filled", key)
	}
	return value
}
