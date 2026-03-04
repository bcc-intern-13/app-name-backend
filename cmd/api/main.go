package main

import (
	bootstrap "github.com/bcc-intern-13/app-name-backend/cmd/bootsrap"
	"github.com/bcc-intern-13/app-name-backend/internal/user/handler"
	"github.com/bcc-intern-13/app-name-backend/internal/user/repository"
	"github.com/bcc-intern-13/app-name-backend/internal/user/service"
)

func main() {
	app := bootstrap.NewApp()
	userRepo := repository.NewUserRepository(app.DB)
	userService := service.NewUserAuthService(userRepo)
	handler.NewAuthHandler(app.Fiber, userService)
	app.Run()
}
