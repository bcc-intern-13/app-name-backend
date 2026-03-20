package bootstrap

import (
	"log"

	"github.com/bcc-intern-13/app-name-backend/config"
	"github.com/bcc-intern-13/app-name-backend/pkg/email"
	"github.com/bcc-intern-13/app-name-backend/pkg/gemini"
	"github.com/bcc-intern-13/app-name-backend/pkg/storage"

	"github.com/bcc-intern-13/app-name-backend/internal/infra/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	Fiber          *fiber.App
	DB             *gorm.DB
	Config         *config.Config
	EmailService   *email.EmailService
	StorageService *storage.StorageService
	geminiService  *gemini.GeminiService
}

func NewApp() *App {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	//migrate and seed
	database.Migrate(db)
	database.Seed(db)

	//email package
	EmailService := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
		cfg.AppURL,
	)

	//storage services package
	storageService := storage.NewStorageService(
		cfg.SupabaseURL,
		cfg.SupabaseServiceRoleKey,
		cfg.StorageBucketCV,
		cfg.StorageBucketAvatar,
	)

	//gemini service package
	geminiService, err := gemini.NewGeminiService(cfg.GeminiAPIKey)
	if err != nil {
		log.Fatal("Failed to initialize Gemini:", err)
	}

	return &App{
		Fiber:          fiber.New(),
		DB:             db,
		Config:         cfg,
		EmailService:   EmailService,
		StorageService: storageService,
		geminiService:  geminiService,
	}
}

func (a *App) Run() {
	log.Printf("Server running on http://localhost:%s", a.Config.Port)
	log.Fatal(a.Fiber.Listen(":" + a.Config.Port))
}
