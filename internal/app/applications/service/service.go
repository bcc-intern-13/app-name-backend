package service

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"

	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/entity"
	jobBoardContract "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/contract"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/bcc-intern-13/app-name-backend/pkg/storage"
	"github.com/google/uuid"
)

type applicationService struct {
	repo           contract.ApplicationRepository
	jobBoardRepo   jobBoardContract.JobBoardRepository
	storageService *storage.StorageService
}

func NewApplicationService(
	repo contract.ApplicationRepository,
	jobBoardRepo jobBoardContract.JobBoardRepository,
	storageService *storage.StorageService,
) contract.ApplicationService {
	return &applicationService{
		repo:           repo,
		jobBoardRepo:   jobBoardRepo,
		storageService: storageService,
	}
}

func (s *applicationService) Submit(userID uuid.UUID, req *dto.SubmitApplicationRequest, cv *multipart.FileHeader) *response.APIError {
	jobID, err := uuid.Parse(req.JobID)
	if err != nil {
		return response.ErrBadRequest("invalid job id")
	}

	job, err := s.jobBoardRepo.FindByID(jobID)
	if err != nil {
		slog.Error("failed to get job", "error", err, "jobID", jobID)
		return response.ErrInternal("failed to get job")
	}
	if job == nil {
		return response.ErrNotFound("job not found")
	}

	existing, err := s.repo.FindByUserIDAndJobID(userID, jobID)
	if err != nil {
		slog.Error("failed to check existing application", "error", err)
		return response.ErrInternal("failed to check existing application")
	}
	if existing != nil {
		return response.ErrConflict("you have already applied to this job")
	}

	// upload cv if provided
	cvURL := ""
	if cv != nil {
		// validate file type (must be pdf files)
		if cv.Header.Get("Content-Type") != "application/pdf" {
			return response.ErrBadRequest("cv must be a PDF file")
		}

		// validate file size (max 5MB)
		if cv.Size > 5*1024*1024 {
			return response.ErrBadRequest("cv file size must be less than 5MB")
		}

		// read file content
		file, err := cv.Open()
		if err != nil {
			slog.Error("failed to open cv file", "error", err)
			return response.ErrInternal("failed to process cv file")
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			slog.Error("failed to read cv file", "error", err)
			return response.ErrInternal("failed to read cv file")
		}

		// upload to the supabase storage
		cvURL, err = s.storageService.UploadCV(
			fmt.Sprintf("%s_%s", userID.String(), jobID.String()),
			fileBytes,
			"application/pdf",
		)

		if err != nil {
			slog.Error("failed to upload cv", "error", err.Error(), "userID", userID) // ← tambah .Error()
			return response.ErrInternal("failed to upload cv")
		}
	}

	application := &entity.Application{
		ID:            uuid.New(),
		UserID:        userID,
		JobID:         jobID,
		CvURL:         cvURL,
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
		result = append(result, dto.ApplicationResponse{
			ID:            app.ID,
			JobID:         app.JobID,
			JobTitle:      app.JobTitle,
			JobCity:       app.JobCity,
			JobType:       app.JobType,
			CompanyName:   app.CompanyName,
			CompanyLogo:   app.CompanyLogo,
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

	if application.UserID != userID {
		return nil, response.ErrUnAuthorized("you are not authorized to view this application")
	}

	return &dto.ApplicationDetailResponse{
		ApplicationResponse: dto.ApplicationResponse{
			ID:            application.ID,
			JobID:         application.JobID,
			JobTitle:      application.JobTitle,
			JobCity:       application.JobCity,
			JobType:       application.JobType,
			CompanyName:   application.CompanyName,
			CompanyLogo:   application.CompanyLogo,
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

	// making sure user can only delete their own application
	if application.UserID != userID {
		return response.ErrUnAuthorized("you are not authorized to delete this application")
	}

	//cancel only allowed for applications with status "Terkirim"
	if application.Status != "Terkirim" {
		return response.ErrBadRequest("can only withdraw application with status Terkirim")
	}

	if err := s.repo.Delete(id); err != nil {
		slog.Error("failed to delete application", "error", err, "id", id)
		return response.ErrInternal("failed to withdraw application")
	}

	return nil
}
