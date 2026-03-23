package dto

import (
	"time"

	"github.com/google/uuid"
)

// CreateOrderResponse → POST /api/payment/create-order
type CreateOrderResponse struct {
	OrderID    uuid.UUID `json:"order_id"`
	PaymentURL string    `json:"payment_url"`
	Amount     int64     `json:"amount"`
}

// OrderResponse → GET /api/payment/history
type OrderResponse struct {
	ID              uuid.UUID `json:"id"`
	MidtransOrderID string    `json:"midtrans_order_id"`
	Amount          int64     `json:"amount"`
	Status          string    `json:"status"`
	PaymentType     string    `json:"payment_type"`
	CreatedAt       time.Time `json:"created_at"`
}
