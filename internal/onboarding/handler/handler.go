package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/onboarding/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

type onboardingHandler struct {
	service dto.OnboardingService
}

func (h *onboardingHandler) submit(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.SubmitOnboardingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if err := h.service.Submit(userID, &req); err != nil {
		switch err.Error() {
		case "onboarding already completed":
			return response.Error(ctx, response.ErrConflict(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	return response.Success(ctx, fiber.StatusCreated, "Onboarding submitted successfully", nil)
}

// to get onboarding answers later
func (h *onboardingHandler) getAnswers(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	profile, err := h.service.GetByUserID(userID)
	if err != nil {
		switch err.Error() {
		case "profile not found":
			return response.Error(ctx, response.ErrNotFound(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	return response.Success(ctx, fiber.StatusOK, "success", profile)
}

func (h *onboardingHandler) update(ctx *fiber.Ctx) error {
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	var req dto.SubmitOnboardingRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if err := h.service.Update(userID, &req); err != nil {
		switch err.Error() {
		case "profile not found":
			return response.Error(ctx, response.ErrNotFound(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	return response.Success(ctx, fiber.StatusOK, "Onboarding updated successfully", nil)
}
