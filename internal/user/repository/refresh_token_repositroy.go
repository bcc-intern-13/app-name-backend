package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) dto.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(token *entity.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*entity.RefreshToken, error) {
	var rt entity.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &rt, err
}

func (r *refreshTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entity.RefreshToken{}).Error
}

func (r *refreshTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&entity.RefreshToken{}).Error
}
