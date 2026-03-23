package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	MidtransOrderID string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"midtrans_order_id"`
	Amount          int64     `gorm:"not null" json:"amount"`
	Status          string    `gorm:"type:varchar(20);default:'PENDING'" json:"status"` // PENDING | PAID | FAILED | EXPIRED
	PaymentURL      string    `gorm:"type:text" json:"payment_url"`
	PaymentType     string    `gorm:"type:varchar(50)" json:"payment_type"` // diisi Midtrans waktu webhook
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
