package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/onboarding/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterOnboardingRoutes(app *fiber.App, service dto.OnboardingService, jwtSecret string) {
	onboarding := app.Group("/api/onboarding", middleware.JWTProtected(jwtSecret))
	h := &onboardingHandler{service: service}

	onboarding.Post("/submit", h.submit)
	onboarding.Get("/answers", h.getAnswers)
	onboarding.Patch("/update", h.update)
}
