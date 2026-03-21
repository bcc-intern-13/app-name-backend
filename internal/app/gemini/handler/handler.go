package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/gemini/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

type cvHandler struct {
	service contract.CVService
}

// POST /api/cv/upload — upload PDF ke storage, TANPA Gemini
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

// POST /api/cv/analyze — panggil Gemini extract dari PDF yang udah diupload
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

// GET /api/cv
func (h *cvHandler) getCV(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetCV(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

// PATCH /api/cv
func (h *cvHandler) updateCV(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.UpdateCVRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}

	result, apiErr := h.service.UpdateCV(ctx.Context(), userID, &req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "CV updated successfully", result)
}

// GET /api/cv/score
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

// GET /api/cv/ai-calls-remaining
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

// POST /api/cv-ai/improve-sentence
func (h *cvHandler) improveSentence(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.ImproveSentenceRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}
	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	result, apiErr := h.service.ImproveSentence(ctx.Context(), userID, &req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

// POST /api/cv-ai/job-match
func (h *cvHandler) jobMatch(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.JobMatchRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}
	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	result, apiErr := h.service.JobMatch(ctx.Context(), userID, &req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

// POST /api/cv-ai/review
func (h *cvHandler) reviewCV(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.ReviewCV(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func getUserID(ctx *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(ctx.Locals("userID").(string))
}
