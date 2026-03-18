package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/job_board/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterJobBoardRoutes(app *fiber.App, service dto.JobBoardService, jwtSecret string) {
	h := &jobBoardHandler{service: service}

	jobBoard := app.Group("/api/job-board", middleware.JWTProtected(jwtSecret))
	jobBoard.Get("/", h.getAll)
	jobBoard.Get("/saved", h.getSavedJobs)
	jobBoard.Get("/:id", h.getByID)
	jobBoard.Post("/:id/save", h.toggleSave)
}
