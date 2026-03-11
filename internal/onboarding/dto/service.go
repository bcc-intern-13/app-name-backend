package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/google/uuid"
)

type OnboardingService interface {
	Submit(userID uuid.UUID, req *SubmitOnboardingRequest) error
	GetByUserID(userID uuid.UUID) (*entity.UserProfile, error)
	Update(userID uuid.UUID, req *SubmitOnboardingRequest) error
}
