package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCareerMappingRoutes(app *fiber.App, service dto.CareerMappingService, jwtSecret string) {
	h := &careerMappingHandler{service: service}

	app.Get("/api/career-mapping/questions", h.getQuestions)

	cm := app.Group("/api/career-mapping", middleware.JWTProtected(jwtSecret))
	cm.Post("/submit", h.submit)
	cm.Get("/result", h.getLatestResult)
}
