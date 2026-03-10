package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"gorm.io/gorm"
)

type verificationTokenRepository struct {
	db *gorm.DB
}

func NewVerificationTokenRepository(db *gorm.DB) dto.VerificationTokenRepository {
	return &verificationTokenRepository{db: db}
}

func (r *verificationTokenRepository) Create(token *entity.VerificationToken) error {
	return r.db.Create(token).Error //note ini kenapa kalo di .error bisa balikin error ?™¡™
}

func (r *verificationTokenRepository) FindByToken(token string) (*entity.VerificationToken, error) {
	var verificationToken entity.VerificationToken
	err := r.db.Where("token = ?", token).First(&verificationToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &verificationToken, err
}

func (r *verificationTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&entity.VerificationToken{}).Error
}
