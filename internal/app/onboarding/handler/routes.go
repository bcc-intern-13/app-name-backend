package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterOnboardingRoutes(app *fiber.App, service contract.OnboardingService, jwtSecret string) {
	onboarding := app.Group("/api/v1/onboarding", middleware.JWTProtected(jwtSecret))
	h := &onboardingHandler{service: service}

	onboarding.Post("/submit", h.submit)
	onboarding.Get("/answers", h.getAnswers)
	onboarding.Patch("/update", h.update)
}
