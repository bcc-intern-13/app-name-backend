package contract

import (
	"context"

	"github.com/bcc-intern-13/app-name-backend/internal/app/payment/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type PaymentService interface {
	// CreateOrder make an order to send it to Xendit,from ther a url will be sent from xendit
	CreateOrder(ctx context.Context, userID uuid.UUID) (*dto.CreateOrderResponse, *response.APIError)

	// HandleWebhook recive notification from xendit and update status order also update to is premium if payment success
	HandleWebhook(ctx context.Context, req *dto.WebhookRequest) *response.APIError

	// GetOrderHistory get order history for user, ordered by created_at desc
	GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]dto.OrderResponse, *response.APIError)
}
