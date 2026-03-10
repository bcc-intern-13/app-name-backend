package bootstrap

import (
	"log"

	"github.com/bcc-intern-13/app-name-backend/config"
	"github.com/bcc-intern-13/app-name-backend/pkg/email"

	"github.com/bcc-intern-13/app-name-backend/internal/infra/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type App struct {
	Fiber        *fiber.App
	DB           *gorm.DB
	Config       *config.Config
	emailservice *email.EmailService
}

func NewApp() *App {
	cfg := config.Load()
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		panic(err)
	}

	database.Migrate(db)

	emailSvc := email.NewEmailService(
		cfg.SMTPHost,
		cfg.SMTPPort,
		cfg.SMTPEmail,
		cfg.SMTPPassword,
		cfg.AppURL,
	)

	return &App{
		Fiber:        fiber.New(),
		DB:           db,
		Config:       cfg,
		emailservice: emailSvc,
	}
}

func (a *App) Run() {
	log.Printf("Server running on http://localhost:%s", a.Config.Port)
	log.Fatal(a.Fiber.Listen(":" + a.Config.Port))
}
