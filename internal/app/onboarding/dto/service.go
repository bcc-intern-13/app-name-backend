package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type OnboardingService interface {
	Submit(userID uuid.UUID, req *SubmitOnboardingRequest) *response.APIError
	GetByUserID(userID uuid.UUID) (*entity.UserProfile, *response.APIError)
	Update(userID uuid.UUID, req *SubmitOnboardingRequest) *response.APIError
}
