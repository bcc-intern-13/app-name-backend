package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/job_board/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/job_board/entity"
	"github.com/google/uuid"
)

type JobBoardRepository interface {
	FindAll(filter dto.JobBoardFilter) ([]entity.JobListing, int64, error)
	FindByID(id uuid.UUID) (*entity.JobListing, error)
	FindCompanyByID(id uuid.UUID) (*entity.Company, error)
	SaveJob(userID, jobID uuid.UUID) error
	UnsaveJob(userID, jobID uuid.UUID) error
	IsJobSaved(userID, jobID uuid.UUID) (bool, error)
	FindSavedJobs(userID uuid.UUID) ([]entity.JobListing, error)
}
