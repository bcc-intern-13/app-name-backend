package contract

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type OnboardingService interface {
	Submit(userID uuid.UUID, req *dto.SubmitOnboardingRequest) *response.APIError
	GetByUserID(userID uuid.UUID) (*entity.UserProfile, *response.APIError)
	Update(userID uuid.UUID, req *dto.SubmitOnboardingRequest) *response.APIError
}
