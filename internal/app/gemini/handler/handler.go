package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type cvHandler struct {
	service contract.CVService
}

// helper to get UUID from token
func getUserID(ctx *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(ctx.Locals("userID").(string))
}

// Uploading cv to gemini without gemini
func (h *cvHandler) uploadCV(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	file, err := ctx.FormFile("cv")
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("cv file is required"), err)
	}

	result, apiErr := h.service.UploadCV(ctx.Context(), userID, file)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusCreated, "CV uploaded successfully", result)
}

// Analyze cv to extract pdf
func (h *cvHandler) analyzeCV(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.AnalyzeCV(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "CV analyzed successfully", result)
}

// Get Score
func (h *cvHandler) getScore(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetScore(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *cvHandler) getAICallsRemaining(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetAICallsRemaining(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *cvHandler) ImproveSentence(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}
	result, apiErr := h.service.ImproveSentence(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}
	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *cvHandler) SuggestKeywords(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}
	result, apiErr := h.service.SuggestKeywords(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}
	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func (h *cvHandler) SummarizeProfile(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}
	result, apiErr := h.service.SummarizeProfile(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}
	return response.Success(ctx, fiber.StatusOK, "success", result)
}
