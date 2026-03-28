package bootstrap

import (
	"log"

	"github.com/bcc-intern-13/WorkAble-backend/config"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/email"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/gemini"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/storage"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/xendit"

	"github.com/bcc-intern-13/WorkAble-backend/internal/infra/database"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type App struct {
	Fiber          *fiber.App
	DB             *gorm.DB
	Config         *config.Config
	EmailService   *email.EmailService
	StorageService *storage.StorageService
	GeminiService  *gemini.GeminiService
	XenditService  *xendit.XenditService
}

func NewApp() *App {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	// migrate and seed
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

	c := cron.New()

	// 0 0 * * means run it every night or 00.00 AM
	_, cronErr := c.AddFunc("0 0 * * *", func() {
		log.Println("⏳ [CRON] Tengah malam tiba! Memulai reset kuota AI Calls...")

		// exec raw sql
		result := db.Exec("UPDATE cvs SET ai_calls_today = 0 WHERE ai_calls_today > 0")

		if result.Error != nil {
			log.Println("❌ [CRON] Duh, gagal ngereset kuota AI:", result.Error)
		} else {
			log.Printf("✅ [CRON] Reset kuota sukses. %d data CV di-reset. User siap gacha AI lagi!\n", result.RowsAffected)
		}
	})

	if cronErr != nil {
		log.Println("⚠️ [CRON] Gagal memasang jam beker:", cronErr)
	} else {
		c.Start() // Nyalakan bekernya di background
		log.Println("✅ [CRON] Jam beker penjaga kuota AI berhasil dipasang di background.")
	}
	// =========================================================================

	return &App{
		Fiber:          fiber.New(),
		DB:             db,
		Config:         cfg,
		EmailService:   EmailService,
		StorageService: storageService,
		GeminiService:  geminiService,
		XenditService:  xenditService,
	}
}

func (a *App) Run() {
	log.Printf("Server running on http://localhost:%s", a.Config.Port)
	log.Fatal(a.Fiber.Listen(":" + a.Config.Port))
}
