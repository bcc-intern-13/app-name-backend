package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

type applicationHandler struct {
	service contract.ApplicationService
}

func (h *applicationHandler) submit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.SubmitApplicationRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	// ambil file CV
	cv, err := ctx.FormFile("cv")
	// if err != nil {
	// 	return response.Error(ctx, response.ErrBadRequest("cv file is required"), err)
	// }

	if apiErr := h.service.Submit(userID, &req, cv); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusCreated, "Application submitted successfully", nil)
}

func (h *applicationHandler) getAll(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	// optional filter by status
	status := ctx.Query("status")

	result, apiErr := h.service.GetAll(userID, status)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *applicationHandler) getByID(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid application id"), err)
	}

	result, apiErr := h.service.GetByID(userID, id)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *applicationHandler) delete(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid application id"), err)
	}

	if apiErr := h.service.Delete(userID, id); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Application withdrawn successfully", nil)
}
