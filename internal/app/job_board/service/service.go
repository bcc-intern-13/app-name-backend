package service

import (
	"log/slog"

	companyContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/company/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type jobBoardService struct {
	repo        contract.JobBoardRepository
	companyRepo companyContract.CompanyRepository // ← tambah

}

func NewJobBoardService(repo contract.JobBoardRepository) contract.JobBoardService {
	return &jobBoardService{repo: repo}
}

func (s *jobBoardService) GetAll(filter dto.JobBoardFilter, userID uuid.UUID) (*dto.PaginatedJobResponse, *response.APIError) {
	jobs, total, err := s.repo.FindAll(filter)
	if err != nil {
		slog.Error("failed to get job listings", "error", err)
		return nil, response.ErrInternal("failed to get job listings")
	}

	var result []dto.JobListingResponse
	for _, job := range jobs {
		result = append(result, dto.JobListingResponse{
			ID:                 job.ID,
			CompanyID:          job.CompanyID,
			CompanyName:        job.CompanyName, // ← langsung dari JOIN
			CompanyLogo:        job.CompanyLogo, // ← langsung dari JOIN
			Title:              job.Title,
			City:               job.City,
			JobType:            job.JobType,
			JobField:           job.JobField,
			Salary:             job.Salary,
			AcceptedDisability: job.AcceptedDisability,
			AccessibilityLabel: job.AccessibilityLabel,
			CreatedAt:          job.CreatedAt,
		})
	}

	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	return &dto.PaginatedJobResponse{
		Data:  result,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *jobBoardService) GetByID(id uuid.UUID) (*dto.JobListingDetailResponse, *response.APIError) {
	job, err := s.repo.FindByID(id)
	if err != nil {
		slog.Error("failed to get job by id", "error", err, "jobID", id)
		return nil, response.ErrInternal("failed to get job")
	}
	if job == nil {
		return nil, response.ErrNotFound("job not found")
	}

	company, _ := s.companyRepo.FindCompanyByID(job.CompanyID)
	companyName := ""
	companyLogo := ""
	if company != nil {
		companyName = company.Name
		companyLogo = company.LogoURL
	}

	return &dto.JobListingDetailResponse{
		JobListingResponse: dto.JobListingResponse{
			ID:                 job.ID,
			CompanyID:          job.CompanyID,
			CompanyName:        companyName,
			CompanyLogo:        companyLogo,
			Title:              job.Title,
			City:               job.City,
			JobType:            job.JobType,
			JobField:           job.JobField,
			Salary:             job.Salary,
			AcceptedDisability: job.AcceptedDisability,
			AccessibilityLabel: job.AccessibilityLabel,
			CreatedAt:          job.CreatedAt,
		},
		Description:   job.Description,
		Qualification: job.Qualification,
	}, nil
}

func (s *jobBoardService) ToggleSave(userID, jobID uuid.UUID) (bool, *response.APIError) {
	job, err := s.repo.FindByID(jobID)
	if err != nil {
		slog.Error("failed to get job", "error", err, "jobID", jobID)
		return false, response.ErrInternal("failed to get job")
	}
	if job == nil {
		return false, response.ErrNotFound("job not found")
	}

	isSaved, err := s.repo.IsJobSaved(userID, jobID)
	if err != nil {
		slog.Error("failed to check saved job", "error", err, "userID", userID, "jobID", jobID)
		return false, response.ErrInternal("failed to check saved job")
	}

	if isSaved {
		if err := s.repo.UnsaveJob(userID, jobID); err != nil {
			slog.Error("failed to unsave job", "error", err, "userID", userID, "jobID", jobID)
			return false, response.ErrInternal("failed to unsave job")
		}
		return false, nil
	}

	if err := s.repo.SaveJob(userID, jobID); err != nil {
		slog.Error("failed to save job", "error", err, "userID", userID, "jobID", jobID)
		return false, response.ErrInternal("failed to save job")
	}
	return true, nil
}

func (s *jobBoardService) GetSavedJobs(userID uuid.UUID) ([]dto.JobListingResponse, *response.APIError) {
	jobs, err := s.repo.FindSavedJobs(userID)
	if err != nil {
		slog.Error("failed to get saved jobs", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get saved jobs")
	}

	var result []dto.JobListingResponse
	for _, job := range jobs {
		company, _ := s.companyRepo.FindCompanyByID(job.CompanyID)
		companyName := ""
		companyLogo := ""
		if company != nil {
			companyName = company.Name
			companyLogo = company.LogoURL
		}

		result = append(result, dto.JobListingResponse{
			ID:                 job.ID,
			CompanyID:          job.CompanyID,
			CompanyName:        companyName,
			CompanyLogo:        companyLogo,
			Title:              job.Title,
			City:               job.City,
			JobType:            job.JobType,
			JobField:           job.JobField,
			Salary:             job.Salary,
			AcceptedDisability: job.AcceptedDisability,
			AccessibilityLabel: job.AccessibilityLabel,
			CreatedAt:          job.CreatedAt,
		})
	}

	return result, nil
}
