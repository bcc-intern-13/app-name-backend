package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/entity"
	"github.com/google/uuid"
)

type OnboardingRepository interface {
	Create(profile *entity.UserProfile) error
	FindByUserID(userID uuid.UUID) (*entity.UserProfile, error)
	Update(profile *entity.UserProfile) error
}
