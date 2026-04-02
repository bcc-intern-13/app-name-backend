package service

import (
	"log/slog"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/dto"
	jobDto "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type companyService struct {
	repo contract.CompanyRepository
}

func NewCompanyService(repo contract.CompanyRepository) contract.CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) GetByID(id uuid.UUID) (*dto.CompanyResponse, *response.APIError) {
	company, err := s.repo.FindCompanyByID(id)
	if err != nil {
		slog.Error("failed to get company", "error", err, "id", id)
		return nil, response.ErrInternal("failed to get company")
	}
	if company == nil {
		return nil, response.ErrNotFound("company not found")
	}

	// Get jobs by company id
	jobs, err := s.repo.FindActiveJobsByCompanyID(id)
	if err != nil {
		slog.Error("failed to get company jobs", "error", err, "id", id)
		return nil, response.ErrInternal("failed to get company jobs")
	}

	// Get "Similar" Companies (get companies except the current company)
	var similarCompanies []dto.CompanyPreviewResponse
	otherCompanies, err := s.repo.FindAllCompaniesExcluding(id)
	if err != nil {
		slog.Error("failed to get other companies", "error", err, "id", id)
	} else {
		for _, c := range otherCompanies {
			similarCompanies = append(similarCompanies, dto.CompanyPreviewResponse{
				ID:       c.ID,
				Name:     c.Name,
				LogoURL:  c.LogoURL,
				Industry: c.Industry,
				Location: c.Location,
			})
		}
	}

	var jobListings []jobDto.JobListingResponse
	for _, job := range jobs {
		jobListings = append(jobListings, jobDto.JobListingResponse{
			ID:                 job.ID,
			CompanyID:          job.CompanyID,
			CompanyName:        job.CompanyName,
			CompanyLogo:        job.CompanyLogo,
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

	return &dto.CompanyResponse{
		ID:                 company.ID,
		Name:               company.Name,
		LogoURL:            company.LogoURL,
		Description:        company.Description,
		Industry:           company.Industry,
		Size:               company.Size,
		Location:           company.Location,
		Website:            company.Website,
		AcceptedDisability: company.AcceptedDisability,
		AccessibilityLabel: company.AccessibilityLabel,
		CreatedAt:          company.CreatedAt,
		JobListings:        jobListings,
		SimilarCompanies:   similarCompanies,
	}, nil
}
