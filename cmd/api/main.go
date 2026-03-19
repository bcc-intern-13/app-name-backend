package main

import (
	bootstrap "github.com/bcc-intern-13/app-name-backend/cmd/bootsrap"

	// user domain packages
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/handler"
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/repository"
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/service"

	// onboarding domain packages
	onboardingHandler "github.com/bcc-intern-13/app-name-backend/internal/app/onboarding/handler"
	onboardingRepository "github.com/bcc-intern-13/app-name-backend/internal/app/onboarding/repository"
	onboardingService "github.com/bcc-intern-13/app-name-backend/internal/app/onboarding/service"

	// career mapping domain packages
	careerMappingHandler "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/handler"
	careerMappingRepository "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/repository"
	careerMappingService "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/service"

	// job board domain packages
	jobBoardHandler "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/handler"
	jobBoardRepository "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/repository"
	jobBoardService "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/service"

	// home domain packages
	homeHandler "github.com/bcc-intern-13/app-name-backend/internal/app/home/handler"
	homeService "github.com/bcc-intern-13/app-name-backend/internal/app/home/service"

	//applications domain packages
	applicationHandler "github.com/bcc-intern-13/app-name-backend/internal/app/applications/handler"
	applicationRepository "github.com/bcc-intern-13/app-name-backend/internal/app/applications/repository"
	applicationService "github.com/bcc-intern-13/app-name-backend/internal/app/applications/service"

	//company domain packages
	companyHandler "github.com/bcc-intern-13/app-name-backend/internal/app/company/handler"
	companyRepository "github.com/bcc-intern-13/app-name-backend/internal/app/company/repository"
	companyService "github.com/bcc-intern-13/app-name-backend/internal/app/company/service"
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

	// onboarding routes
	handler.RegisterRoutes(app.Fiber, userService, app.Config.JWTSecret)
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
	homeSvc := homeService.NewHomeService(onboardingRepo, jobBoardSvc, careerMappingSvc)

	// home routes
	homeHandler.RegisterHomeRoutes(app.Fiber, homeSvc, app.Config.JWTSecret)

	//applications domain
	applicationRepo := applicationRepository.NewApplicationRepository(app.DB)
	applicationSvc := applicationService.NewApplicationService(applicationRepo, jobBoardRepo)

	//applications routes
	applicationHandler.RegisterApplicationRoutes(app.Fiber, applicationSvc, app.Config.JWTSecret)

	//company doomain
	companyRepo := companyRepository.NewCompanyRepository(app.DB)
	companySvc := companyService.NewCompanyService(companyRepo)

	// company routes
	companyHandler.RegisterCompanyRoutes(app.Fiber, companySvc, app.Config.JWTSecret)

	app.Run()
}
