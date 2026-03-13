package dto

import "github.com/google/uuid"

type JobBoardService interface {
	GetAll(filter JobBoardFilter, userID uuid.UUID) (*PaginatedJobResponse, error)
	GetByID(id uuid.UUID) (*JobListingDetailResponse, error)
	ToggleSave(userID, jobID uuid.UUID) (bool, error)
	GetSavedJobs(userID uuid.UUID) ([]JobListingResponse, error)
}
