package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterApplicationRoutes(app *fiber.App, service contract.ApplicationService, jwtSecret string) {
	h := &applicationHandler{service: service}

	applications := app.Group("/api/v1/applications", middleware.JWTProtected(jwtSecret))
	applications.Post("/", h.submit)
	applications.Get("/", h.getAll)
	applications.Get("/:id", h.getByID)
	applications.Delete("/:id", h.delete)
}
