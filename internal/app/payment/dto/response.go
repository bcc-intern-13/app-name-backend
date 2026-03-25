package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderResponse struct {
	OrderID    uuid.UUID `json:"order_id"`
	PaymentURL string    `json:"payment_url"`
	Amount     int64     `json:"amount"`
}
type OrderResponse struct {
	ID               uuid.UUID `json:"id"`
	XenditExternalID string    `json:"xendit_external_id"`

	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	PaymentType string    `json:"payment_type"`
	CreatedAt   time.Time `json:"created_at"`
}
