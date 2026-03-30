package main

import (
	"context"
	"log/slog"

	bootstrap "github.com/bcc-intern-13/WorkAble-backend/cmd/bootsrap"
	"github.com/robfig/cron/v3"

	// user domain packages
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/handler"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/repository"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/service"

	// onboarding domain packages
	onboardingHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/handler"
	onboardingRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/repository"
	onboardingService "github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/service"

	// career mapping domain packages
	careerMappingHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/handler"
	careerMappingRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/repository"
	careerMappingService "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/service"

	// job board domain packages
	jobBoardHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/handler"
	jobBoardRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/repository"
	jobBoardService "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/service"

	// home domain packages
	homeHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/home/handler"
	homeService "github.com/bcc-intern-13/WorkAble-backend/internal/app/home/service"

	//applications domain packages
	applicationHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/handler"
	applicationRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/repository"
	applicationService "github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/service"

	//company domain packages
	companyHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/company/handler"
	companyRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/company/repository"
	companyService "github.com/bcc-intern-13/WorkAble-backend/internal/app/company/service"

	//smartprofile domain packages
	smartProfileHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/smart_profile/handler"
	smartProfileService "github.com/bcc-intern-13/WorkAble-backend/internal/app/smart_profile/service"

	//company domain packages
	geminiHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/handler"
	geminiRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/repository"
	geminiService "github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/service"

	//payment domain packages
	paymentHandler "github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/handler"
	paymentRepository "github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/repository"
	paymentService "github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/service"

	//allow cors for fe implementaion
	"github.com/gofiber/fiber/v2/middleware/cors"
	//google oatuh packages
)

func main() {
	app := bootstrap.NewApp()

	// CORS configuration
	app.Fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, ngrok-skip-browser-warning",
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	}))

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
		app.StorageService,
	)

	// onboarding domain
	onboardingRepo := onboardingRepository.NewOnboardingRepository(app.DB)
	onboardingSvc := onboardingService.NewOnboardingService(onboardingRepo, userRepo)

	// onboarding routes
	handler.RegisterRoutes(app.Fiber, userService, app.Config.JWTSecret, app.GoogleOAuthService)
	onboardingHandler.RegisterOnboardingRoutes(app.Fiber, onboardingSvc, app.Config.JWTSecret)

	// career mapping domain
	careerMappingRepo := careerMappingRepository.NewCareerMappingRepository(app.DB)
	careerMappingSvc := careerMappingService.NewCareerMappingService(careerMappingRepo)

	// career mapping routes
	careerMappingHandler.RegisterCareerMappingRoutes(app.Fiber, careerMappingSvc, app.Config.JWTSecret)

	// job board domain
	jobBoardRepo := jobBoardRepository.NewJobBoardRepository(app.DB)
	jobBoardSvc := jobBoardService.NewJobBoardService(jobBoardRepo)

	// job board routes
	jobBoardHandler.RegisterJobBoardRoutes(app.Fiber, jobBoardSvc, app.Config.JWTSecret)

	// home domain
	homeSvc := homeService.NewHomeService(onboardingRepo, jobBoardSvc, careerMappingSvc, userRepo)

	// home routes
	homeHandler.RegisterHomeRoutes(app.Fiber, homeSvc, app.Config.JWTSecret)

	//applications domain
	applicationRepo := applicationRepository.NewApplicationRepository(app.DB)
	applicationSvc := applicationService.NewApplicationService(applicationRepo, jobBoardRepo, app.StorageService)

	//applications routes
	applicationHandler.RegisterApplicationRoutes(app.Fiber, applicationSvc, app.Config.JWTSecret)

	//company doomain
	companyRepo := companyRepository.NewCompanyRepository(app.DB)
	companySvc := companyService.NewCompanyService(companyRepo)

	// company routes
	companyHandler.RegisterCompanyRoutes(app.Fiber, companySvc, app.Config.JWTSecret)

	// smart profile domain
	smartProfileSvc := smartProfileService.NewSmartProfileService(onboardingRepo, careerMappingSvc, userRepo)

	// smart-profile routes
	smartProfileHandler.RegisterSmartProfileRoutes(app.Fiber, smartProfileSvc, app.Config.JWTSecret)

	//gemini domain
	geminiRepo := geminiRepository.NewCVRepository(app.DB)
	geminiService := geminiService.NewCVService(geminiRepo, app.GeminiService, app.StorageService, userRepo)

	c := cron.New()

	_, cronErr := c.AddFunc("0 0 * * *", func() {
		slog.Info("00.00 Time to reset the AI call QUota..")

		// use the method from repository
		errReset := geminiRepo.ResetAICalls(context.Background())

		if errReset != nil {
			slog.Error("failed to reset AI calls", "error", errReset)
		} else {
			slog.Info("AI Call Quota have been reset to 0 Succes")
		}
	})

	if cronErr != nil {
		slog.Error("failed to setup waker time", "error", cronErr)
	} else {
		c.Start()
		slog.Info("Time waker have been succesfully initialized")
	}

	// gemini routes
	geminiHandler.RegisterRoutes(app.Fiber, geminiService, app.Config.JWTSecret, userRepo)

	//payment domain
	orderRepo := paymentRepository.NewOrderRepository(app.DB)
	paymentService := paymentService.NewPaymentService(orderRepo, userRepo, app.XenditService, app.Config.XenditWebhookToken)

	//payment routes
	paymentHandler.RegisterPaymentRoutes(app.Fiber, paymentService, app.Config.JWTSecret)

	app.Run()
}
