package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/app/onboarding/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type onboardingRepository struct {
	db *gorm.DB
}

func NewOnboardingRepository(db *gorm.DB) contract.OnboardingRepository {
	return &onboardingRepository{db: db}
}

func (r *onboardingRepository) Create(profile *entity.UserProfile) error {
	return r.db.Create(profile).Error
}

func (r *onboardingRepository) FindByUserID(userID uuid.UUID) (*entity.UserProfile, error) {
	var user entity.UserProfile
	err := r.db.Where("user_id = ?", userID).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *onboardingRepository) Update(profile *entity.UserProfile) error {
	return r.db.Save(profile).Error
}
