package contract

import (
	"context"

	"github.com/bcc-intern-13/app-name-backend/internal/app/payment/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type PaymentService interface {
	// CreateOrder → buat order baru, kirim ke Midtrans, return payment_url
	CreateOrder(ctx context.Context, userID uuid.UUID) (*dto.CreateOrderResponse, *response.APIError)

	// HandleWebhook → terima notifikasi dari Midtrans, update status order + is_premium user
	HandleWebhook(ctx context.Context, req *dto.WebhookRequest) *response.APIError

	// GetOrderHistory → riwayat order user
	GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]dto.OrderResponse, *response.APIError)
}
