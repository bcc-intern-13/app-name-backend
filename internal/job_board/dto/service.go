package dto

import (
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type JobBoardService interface {
	GetAll(filter JobBoardFilter, userID uuid.UUID) (*PaginatedJobResponse, *response.APIError)
	GetByID(id uuid.UUID) (*JobListingDetailResponse, *response.APIError)
	ToggleSave(userID, jobID uuid.UUID) (bool, *response.APIError)
	GetSavedJobs(userID uuid.UUID) ([]JobListingResponse, *response.APIError)
}
