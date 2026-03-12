package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/entity"
	"github.com/google/uuid"
)

type CareerMappingService interface {
	GetQuestions() ([]entity.CareerMappingQuestion, error)
	Submit(userID uuid.UUID, req *SubmitCareerMappingRequest) (*CareerMappingResponse, error)
	GetLatestResult(userID uuid.UUID) (*CareerMappingResponse, error)
}
