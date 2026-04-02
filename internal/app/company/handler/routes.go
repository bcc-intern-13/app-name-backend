package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterCompanyRoutes(app *fiber.App, service contract.CompanyService, jwtSecret string) {
	h := &companyHandler{service: service}

	company := app.Group("/api/v1/companies", middleware.JWTProtected(jwtSecret))
	company.Get("/:id", h.getByID)
}
