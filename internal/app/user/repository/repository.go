package repository

import (
	"errors"
	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.UserRepository {
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

// repository for update is user premium
func (r *userRepository) UpdateIsPremium(userID uuid.UUID) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Update("is_premium", true).Error
}

// repository for updating premoum status if  premium has already expired
func (r *userRepository) UpdatePremiumStatus(userID uuid.UUID, isPremium bool, expiresAt *time.Time) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_premium":         isPremium,
			"premium_expires_at": expiresAt,
		}).Error
}

// update function for google auth
func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// find reset paswword repository method
func (r *userRepository) FindByResetToken(token string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
