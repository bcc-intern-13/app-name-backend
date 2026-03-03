package bootstrap

// import (
// 	"log"

// 	"app-name/internal/domain/config"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/yourname/yourproject/infra/database"
// 	"gorm.io/gorm"
// )

// type App struct {
// 	Fiber  *fiber.App
// 	DB     *gorm.DB
// 	Config *config.Config
// }

// func NewApp() *App {
// 	cfg := config.Load()
// 	db := database.Connect(cfg.DatabaseURL)
// 	database.Migrate(db)

// 	return &App{
// 		Fiber:  fiber.New(),
// 		DB:     db,
// 		Config: cfg,
// 	}
// }

// func (a *App) Run() {
// 	log.Printf("🚀 Server running on http://localhost:%s", a.Config.Port)
// 	log.Fatal(a.Fiber.Listen(":" + a.Config.Port))
// }
