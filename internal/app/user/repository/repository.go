package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/app/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/entity"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) dto.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error //use json field id
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// token update status verified true
func (r *userRepository) UpdateVerified(userId uuid.UUID) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userId).
		Update("is_verified", true).Error
}

// upadte onboarding completed true
func (r *userRepository) UpdateOnboardingCompleted(userID uuid.UUID) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("onboarding_completed", true).Error
}
