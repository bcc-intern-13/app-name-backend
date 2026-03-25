package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/payment/entity"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *entity.Order) error
	FindByID(id uuid.UUID) (*entity.Order, error)
	FindByXenditExternalID(xenditExternalID string) (*entity.Order, error)
	FindByUserID(userID uuid.UUID) ([]entity.Order, error)
	Update(order *entity.Order) error
	FindPendingByUserID(userID uuid.UUID) (*entity.Order, error)
}
