package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/smart_profile/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type SmartProfileService interface {
	GetByUserID(userID uuid.UUID) (*dto.SmartProfileResponse, *response.APIError)
}
