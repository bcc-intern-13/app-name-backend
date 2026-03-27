package contract

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type JobBoardService interface {
	GetAll(filter dto.JobBoardFilter, userID uuid.UUID) (*dto.PaginatedJobResponse, *response.APIError)
	GetByID(id uuid.UUID) (*dto.JobListingDetailResponse, *response.APIError)
	ToggleSave(userID, jobID uuid.UUID) (bool, *response.APIError)
	GetSavedJobs(userID uuid.UUID) ([]dto.JobListingResponse, *response.APIError)
}
