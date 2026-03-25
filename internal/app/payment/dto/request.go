package dto

type WebhookRequest struct {
	ID            string  `json:"id"`
	ExternalID    string  `json:"external_id"`
	Status        string  `json:"status"` // PAID or EXPIRED
	PaymentMethod string  `json:"payment_method"`
	Amount        float64 `json:"amount"`
	CallbackToken string  `json:"-"` // take from hader x-callback-token
}
