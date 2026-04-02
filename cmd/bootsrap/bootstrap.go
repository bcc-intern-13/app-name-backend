package bootstrap

import (
	"log"
	"log/slog"

	"github.com/bcc-intern-13/WorkAble-backend/config"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/email"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/gemini"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/oauth"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/storage"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/xendit"

	"github.com/bcc-intern-13/WorkAble-backend/internal/infra/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	Fiber              *fiber.App
	DB                 *gorm.DB
	Config             *config.Config
	EmailService       *email.EmailService
	StorageService     *storage.StorageService
	GeminiService      *gemini.GeminiService
	XenditService      *xendit.XenditService
	GoogleOAuthService *oauth.GoogleOAuthService
}

func NewApp() *App {
	slog.Info("Loading configs and ENV")
	cfg := config.Load()

	slog.Info("connecting to database")
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// migrate and seed
	slog.Info("Migrating to database")
	database.Migrate(db)
	database.Seed(db)

	// email package
	EmailService := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
		cfg.AppURL,
	)

	// storage services package
	storageService := storage.NewStorageService(
		cfg.SupabaseURL,
		cfg.SupabaseServiceRoleKey,
		cfg.StorageBucketCV,
		cfg.StorageBucketAvatar,
	)

	// gemini service package
	geminiService, err := gemini.NewGeminiService(cfg.GeminiAPIKey)
	// log.Printf("Gemini API Key loaded: %s...", cfg.GeminiAPIKey[:10])
	// if err != nil {
	//  log.Fatal("Failed to initialize Gemini:", err)
	// }

	// xendit service package
	xenditService := xendit.NewXenditService(cfg.XenditSecretKey)

	//google oauth service
	googleOAuth := oauth.NewGoogleOAuthService(
		cfg.GoogleClientID,
		cfg.GoogleClientSecret,
		cfg.GoogleRedirectURL,
	)

	return &App{
		Fiber:              fiber.New(),
		DB:                 db,
		Config:             cfg,
		EmailService:       EmailService,
		StorageService:     storageService,
		GeminiService:      geminiService,
		XenditService:      xenditService,
		GoogleOAuthService: googleOAuth,
	}
}

func (a *App) Run() {
	log.Printf("Server running on http://localhost:%s", a.Config.Port)
	log.Fatal(a.Fiber.Listen(":" + a.Config.Port))
}
