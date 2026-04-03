package repository

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) contract.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) FindByID(id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Where("id = ?", id).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *orderRepository) FindByXenditExternalID(xenditExternalID string) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Where("xendit_external_id = ?", xenditExternalID).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *orderRepository) FindByUserID(userID uuid.UUID) ([]entity.Order, error) {
	var orders []entity.Order
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) FindPendingByUserID(userID uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Where("user_id = ? AND status = ?", userID, "PENDING").First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &order, err
}

func (r *orderRepository) FinalizePayment(order *entity.Order, expiresAt *time.Time) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		if err := tx.Save(order).Error; err != nil {
			slog.Error("database error: failed to update order status", "error", err, "orderID", order.ID)
			return err
		}

		if order.Status == "PAID" {

			result := tx.Table("users").Where("id = ?", order.UserID).Updates(map[string]interface{}{
				"is_premium":         true,
				"premium_expires_at": expiresAt,
				"updated_at":         time.Now(),
			})

			if result.Error != nil {
				slog.Error("database error: failed to upgrade user to premium", "error", result.Error, "userID", order.UserID)
				return result.Error
			}

			if result.RowsAffected == 0 {
				slog.Warn("database warning: user not found during premium upgrade", "userID", order.UserID)
				return fmt.Errorf("user not found: %s", order.UserID)
			}
		}

		return nil
	})
}
