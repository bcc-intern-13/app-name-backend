package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/contract"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type companyHandler struct {
	service contract.CompanyService
}

func (h *companyHandler) getByID(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid company id"), err)
	}

	result, apiErr := h.service.GetByID(id)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}
