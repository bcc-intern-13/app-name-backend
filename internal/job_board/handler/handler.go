package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/job_board/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type jobBoardHandler struct {
	service dto.JobBoardService
}

func (h *jobBoardHandler) getAll(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var filter dto.JobBoardFilter
	if err := ctx.QueryParser(&filter); err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid query params"), err)
	}

	result, err := h.service.GetAll(filter, userID)
	if err != nil {
		return response.Error(ctx, response.ErrInternal(err.Error()), err)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *jobBoardHandler) getByID(ctx *fiber.Ctx) error {
	jobID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid job id"), err)
	}

	result, err := h.service.GetByID(jobID)
	if err != nil {
		switch err.Error() {
		case "job not found":
			return response.Error(ctx, response.ErrNotFound(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *jobBoardHandler) toggleSave(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	jobID, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid job id"), err)
	}

	isSaved, err := h.service.ToggleSave(userID, jobID)
	if err != nil {
		switch err.Error() {
		case "job not found":
			return response.Error(ctx, response.ErrNotFound(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	msg := "job saved"
	if !isSaved {
		msg = "job unsaved"
	}

	return response.Success(ctx, fiber.StatusOK, msg, nil)
}

func (h *jobBoardHandler) getSavedJobs(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, err := h.service.GetSavedJobs(userID)
	if err != nil {
		return response.Error(ctx, response.ErrInternal(err.Error()), err)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}
