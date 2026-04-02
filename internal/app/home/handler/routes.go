package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/home/service"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterHomeRoutes(app *fiber.App, svc service.HomeService, jwtSecret string) {
	h := &homeHandler{service: svc}

	home := app.Group("/api/v1/home", middleware.JWTProtected(jwtSecret))
	home.Get("/summary", h.getSummary)
}
