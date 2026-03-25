package entity

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID           uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
	XenditExternalID string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"xendit_external_id"`

	Amount int64  `gorm:"not null" json:"amount"`
	Status string `gorm:"type:varchar(20);default:'PENDING'"` // PENDING or PAID or EXPIRED or FAILED

	PaymentURL  string    `gorm:"type:text" json:"payment_url"`
	PaymentType string    `gorm:"type:varchar(50)" json:"payment_type"` // xendit will fill it from webhook endpoint
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
