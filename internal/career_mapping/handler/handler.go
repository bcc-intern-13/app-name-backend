package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

type careerMappingHandler struct {
	service dto.CareerMappingService
}

func (h *careerMappingHandler) getQuestions(ctx *fiber.Ctx) error {
	questions, err := h.service.GetQuestions()
	if err != nil {
		return response.Error(ctx, response.ErrInternal("failed to get questions"), err)
	}
	return response.Success(ctx, fiber.StatusOK, "success", questions)
}

func (h *careerMappingHandler) submit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.SubmitCareerMappingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	result, err := h.service.Submit(userID, &req)
	if err != nil {
		return response.Error(ctx, response.ErrInternal(err.Error()), err)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *careerMappingHandler) getLatestResult(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, err := h.service.GetLatestResult(userID)
	if err != nil {
		switch err.Error() {
		case "result not found":
			return response.Error(ctx, response.ErrNotFound(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}
