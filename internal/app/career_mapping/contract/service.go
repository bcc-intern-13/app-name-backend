package contract

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/entity"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type CareerMappingService interface {
	GetQuestions() ([]entity.CareerMappingQuestion, *response.APIError)
	Submit(userID uuid.UUID, req *dto.SubmitCareerMappingRequest) (*dto.CareerMappingResponse, *response.APIError)
	GetLatestResult(userID uuid.UUID) (*dto.CareerMappingResponse, *response.APIError)
	GetHistory(userID uuid.UUID) ([]dto.CareerMappingResponse, *response.APIError)
}
