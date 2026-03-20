package bootstrap

import (
	"log"

	"github.com/bcc-intern-13/app-name-backend/config"
	"github.com/bcc-intern-13/app-name-backend/pkg/email"
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
	EmailSvc := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
		cfg.AppURL,
	)

	storageService := storage.NewStorageService(
		cfg.DatabaseURL,
		cfg.SupabaseServiceRoleKey,
		cfg.StorageBucketCV,
		cfg.StorageBucketAvatar,
	)

	return &App{
		Fiber:          fiber.New(),
		DB:             db,
		Config:         cfg,
		EmailService:   EmailSvc,
		StorageService: storageService,
	}
}

func (a *App) Run() {
	log.Printf("Server running on http://localhost:%s", a.Config.Port)
	log.Fatal(a.Fiber.Listen(":" + a.Config.Port))
}
