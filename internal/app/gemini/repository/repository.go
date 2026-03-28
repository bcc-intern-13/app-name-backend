package repository

import (
	"context"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/gemini/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type cvRepository struct {
	db *gorm.DB
}

func NewCVRepository(db *gorm.DB) contract.CVRepository {
	return &cvRepository{db: db}
}

func (r *cvRepository) FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.CV, error) {
	var cv entity.CV
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cv).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cv, nil
}

func (r *cvRepository) Create(ctx context.Context, cv *entity.CV) error {
	return r.db.WithContext(ctx).Create(cv).Error
}

func (r *cvRepository) Update(ctx context.Context, cv *entity.CV) error {
	return r.db.WithContext(ctx).Save(cv).Error
}

func (r *cvRepository) ResetAICalls(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&entity.CV{}).
		Where("ai_calls_today > ?", 0).
		Update("ai_calls_today", 0).Error
}
