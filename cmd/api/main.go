package main

import (
	bootstrap "github.com/bcc-intern-13/app-name-backend/cmd/bootsrap"

	// user domain packages
	"github.com/bcc-intern-13/app-name-backend/internal/user/handler"
	"github.com/bcc-intern-13/app-name-backend/internal/user/repository"
	"github.com/bcc-intern-13/app-name-backend/internal/user/service"

	// onboarding domain packages
	onboardingHandler "github.com/bcc-intern-13/app-name-backend/internal/onboarding/handler"
	onboardingRepository "github.com/bcc-intern-13/app-name-backend/internal/onboarding/repository"
	onboardingService "github.com/bcc-intern-13/app-name-backend/internal/onboarding/service"
)

func main() {
	app := bootstrap.NewApp()

	//user domain
	userRepo := repository.NewUserRepository(app.DB)
	refreshTokenRepo := repository.NewRefreshTokenRepository(app.DB)
	verificationTokenRepo := repository.NewVerificationTokenRepository(app.DB)
	userService := service.NewUserAuthService(
		userRepo,
		app.Config.JWTSecret,
		refreshTokenRepo,
		verificationTokenRepo,
		app.EmailService,
	)

	// onboarding domain
	onboardingRepo := onboardingRepository.NewOnboardingRepository(app.DB)
	onboardingSvc := onboardingService.NewOnboardingService(onboardingRepo, userRepo)

	// routes
	handler.RegisterRoutes(app.Fiber, userService, app.Config.JWTSecret)
	onboardingHandler.RegisterOnboardingRoutes(app.Fiber, onboardingSvc, app.Config.JWTSecret) // ← tambah

	app.Run()
}
