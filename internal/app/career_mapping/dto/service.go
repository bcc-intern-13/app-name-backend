package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type CareerMappingService interface {
	GetQuestions() ([]entity.CareerMappingQuestion, *response.APIError)
	Submit(userID uuid.UUID, req *SubmitCareerMappingRequest) (*CareerMappingResponse, *response.APIError)
	GetLatestResult(userID uuid.UUID) (*CareerMappingResponse, *response.APIError)
}
