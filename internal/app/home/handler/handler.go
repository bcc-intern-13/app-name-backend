package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/home/service"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type homeHandler struct {
	service service.HomeService
}

func (h *homeHandler) getSummaryHome(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetSummary(userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}
