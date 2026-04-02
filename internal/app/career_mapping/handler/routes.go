package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCareerMappingRoutes(app *fiber.App, service contract.CareerMappingService, jwtSecret string) {
	h := &careerMappingHandler{service: service}

	app.Get("/api/v1/career-mapping/questions", h.getQuestions)

	cm := app.Group("/api/career-mapping", middleware.JWTProtected(jwtSecret))
	cm.Post("/submit", h.submit)
	cm.Get("/result", h.getLatestResult)
	cm.Get("/history", h.getHistory)
}
