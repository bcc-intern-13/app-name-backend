package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

type careerMappingHandler struct {
	service contract.CareerMappingService
}

func (h *careerMappingHandler) getQuestions(ctx *fiber.Ctx) error {
	questions, apiErr := h.service.GetQuestions()
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
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

	result, apiErr := h.service.Submit(userID, &req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *careerMappingHandler) getLatestResult(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetLatestResult(userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *careerMappingHandler) getHistory(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetHistory(userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}
