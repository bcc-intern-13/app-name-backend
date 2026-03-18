package service

import (
	"log/slog"
	"mime/multipart"

	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/entity"
	jobBoardContract "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/contract"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type applicationService struct {
	repo         contract.ApplicationRepository
	jobBoardRepo jobBoardContract.JobBoardRepository
}

func NewApplicationService(
	repo contract.ApplicationRepository,
	jobBoardRepo jobBoardContract.JobBoardRepository,
) contract.ApplicationService {
	return &applicationService{
		repo:         repo,
		jobBoardRepo: jobBoardRepo,
	}
}

func (s *applicationService) Submit(userID uuid.UUID, req *dto.SubmitApplicationRequest, cv *multipart.FileHeader) *response.APIError {
	// parse job id
	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		return response.ErrBadRequest("invalid job id")
	}

	// cek job exists
	job, err := s.jobBoardRepo.FindByID(jobID)
	if err != nil {
		slog.Error("failed to get job", "error", err, "jobID", jobID)
		return response.ErrInternal("failed to get job")
	}
	if job == nil {
		return response.ErrNotFound("job not found")
	}

	// cek sudah pernah lamar
	existing, err := s.repo.FindByUserIDAndJobID(userID, jobID)
	if err != nil {
		slog.Error("failed to check existing application", "error", err)
		return response.ErrInternal("failed to check existing application")
	}
	if existing != nil {
		return response.ErrConflict("you have already applied to this job")
	}

	// untuk sekarang cv_url dikosongkan dulu, nanti diisi setelah storage siap
	application := &entity.Application{
		ID:            uuid.New(),
		UserID:        userID,
		JobID:         jobID,
		CvURL:         "",
		PortfolioLink: req.PortfolioLink,
		Status:        "Terkirim",
	}

	if err := s.repo.Create(application); err != nil {
		slog.Error("failed to create application", "error", err, "userID", userID)
		return response.ErrInternal("failed to submit application")
	}

	return nil
}

func (s *applicationService) GetAll(userID uuid.UUID, status string) ([]dto.ApplicationResponse, *response.APIError) {
	applications, err := s.repo.FindAllByUserID(userID, status)
	if err != nil {
		slog.Error("failed to get applications", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get applications")
	}

	var result []dto.ApplicationResponse
	for _, app := range applications {
		// ambil data job
		jobTitle := ""
		companyName := ""
		companyLogo := ""

		job, err := s.jobBoardRepo.FindByID(app.JobID)
		if err != nil {
			slog.Error("failed to get job for application", "error", err, "jobID", app.JobID)
		} else if job != nil {
			jobTitle = job.Title
			company, err := s.jobBoardRepo.FindCompanyByID(job.CompanyID)
			if err != nil {
				slog.Error("failed to get company", "error", err, "companyID", job.CompanyID)
			} else if company != nil {
				companyName = company.Name
				companyLogo = company.LogoURL
			}
		}

		result = append(result, dto.ApplicationResponse{
			ID:            app.ID,
			JobID:         app.JobID,
			JobTitle:      jobTitle,
			CompanyName:   companyName,
			CompanyLogo:   companyLogo,
			Status:        app.Status,
			PortfolioLink: app.PortfolioLink,
			CvURL:         app.CvURL,
			CreatedAt:     app.CreatedAt,
			UpdatedAt:     app.UpdatedAt,
		})
	}

	return result, nil
}

func (s *applicationService) GetByID(userID, id uuid.UUID) (*dto.ApplicationDetailResponse, *response.APIError) {
	application, err := s.repo.FindByID(id)
	if err != nil {
		slog.Error("failed to get application", "error", err, "id", id)
		return nil, response.ErrInternal("failed to get application")
	}
	if application == nil {
		return nil, response.ErrNotFound("application not found")
	}

	// pastikan application milik user ini
	if application.UserID != userID {
		return nil, response.ErrUnAuthorized("you are not authorized to view this application")
	}

	// ambil data job
	jobTitle := ""
	companyName := ""
	companyLogo := ""

	job, err := s.jobBoardRepo.FindByID(application.JobID)
	if err != nil {
		slog.Error("failed to get job", "error", err, "jobID", application.JobID)
	} else if job != nil {
		jobTitle = job.Title
		company, err := s.jobBoardRepo.FindCompanyByID(job.CompanyID)
		if err != nil {
			slog.Error("failed to get company", "error", err, "companyID", job.CompanyID)
		} else if company != nil {
			companyName = company.Name
			companyLogo = company.LogoURL
		}
	}

	return &dto.ApplicationDetailResponse{
		ApplicationResponse: dto.ApplicationResponse{
			ID:            application.ID,
			JobID:         application.JobID,
			JobTitle:      jobTitle,
			CompanyName:   companyName,
			CompanyLogo:   companyLogo,
			Status:        application.Status,
			PortfolioLink: application.PortfolioLink,
			CvURL:         application.CvURL,
			CreatedAt:     application.CreatedAt,
			UpdatedAt:     application.UpdatedAt,
		},
		InterviewLink:   application.InterviewLink,
		InterviewPdfURL: application.InterviewPdfURL,
		InterviewNotes:  application.InterviewNotes,
	}, nil
}

func (s *applicationService) Delete(userID, id uuid.UUID) *response.APIError {
	application, err := s.repo.FindByID(id)
	if err != nil {
		slog.Error("failed to get application", "error", err, "id", id)
		return response.ErrInternal("failed to get application")
	}
	if application == nil {
		return response.ErrNotFound("application not found")
	}

	// pastikan application milik user ini
	if application.UserID != userID {
		return response.ErrUnAuthorized("you are not authorized to delete this application")
	}

	// hanya bisa tarik kalau status masih Terkirim
	if application.Status != "Terkirim" {
		return response.ErrBadRequest("can only withdraw application with status Terkirim")
	}

	if err := s.repo.Delete(id); err != nil {
		slog.Error("failed to delete application", "error", err, "id", id)
		return response.ErrInternal("failed to withdraw application")
	}

	return nil
}
