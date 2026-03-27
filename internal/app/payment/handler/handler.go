package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type paymentHandler struct {
	service contract.PaymentService
}

// create order
func (h *paymentHandler) createOrder(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.CreateOrder(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusCreated, "order created successfully", result)
}

// handleWebhook recive notification from xendit and update status order also update to is premium if payment success
func (h *paymentHandler) handleWebhook(ctx *fiber.Ctx) error {
	var req dto.WebhookRequest
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid webhook payload"), err)
	}

	//note callback token
	req.CallbackToken = ctx.Get("x-callback-token")

	if apiErr := h.service.HandleWebhook(ctx.Context(), &req); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "webhook processed", nil)
}

// history
func (h *paymentHandler) getOrderHistory(ctx *fiber.Ctx) error {
	userID, err := getUserID(ctx)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
	}

	result, apiErr := h.service.GetOrderHistory(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "success", result)
}

func getUserID(ctx *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(ctx.Locals("userID").(string))
}
