package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/smart_profile/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterSmartProfileRoutes(app *fiber.App, service contract.SmartProfileService, jwtSecret string) {
	h := &smartProfileHandler{service: service}

	smartProfile := app.Group("/api/smart-profile", middleware.JWTProtected(jwtSecret))
	smartProfile.Get("/", h.getSmartProfile)
}
