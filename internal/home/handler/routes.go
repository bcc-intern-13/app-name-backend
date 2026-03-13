package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/home/service"
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterHomeRoutes(app *fiber.App, svc service.HomeService, jwtSecret string) {
	h := &homeHandler{service: svc}

	home := app.Group("/api/home", middleware.JWTProtected(jwtSecret))
	home.Get("/summary", h.getSummary)
}
