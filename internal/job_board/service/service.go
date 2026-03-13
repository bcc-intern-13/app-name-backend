package service

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/job_board/dto"
	"github.com/google/uuid"
)

type jobBoardService struct {
	repo dto.JobBoardRepository
}

func NewJobBoardService(repo dto.JobBoardRepository) dto.JobBoardService {
	return &jobBoardService{repo: repo}
}

func (s *jobBoardService) GetAll(filter dto.JobBoardFilter, userID uuid.UUID) (*dto.PaginatedJobResponse, error) {
	jobs, total, err := s.repo.FindAll(filter)
	if err != nil {
		return nil, errors.New("failed to get job listings")
	}

	var result []dto.JobListingResponse
	for _, job := range jobs {
		company, _ := s.repo.FindCompanyByID(job.CompanyID)
		companyName := ""
		companyLogo := ""
		if company != nil {
			companyName = company.Nama
			companyLogo = company.LogoURL
		}

		result = append(result, dto.JobListingResponse{
			ID:                  job.ID,
			CompanyID:           job.CompanyID,
			CompanyName:         companyName,
			CompanyLogo:         companyLogo,
			Judul:               job.Judul,
			Kota:                job.Kota,
			TipePekerjaan:       job.TipePekerjaan,
			BidangKerja:         job.BidangKerja,
			Gaji:                job.Gaji,
			DisabilitasDiterima: job.DisabilitasDiterima,
			LabelAksesibilitas:  job.LabelAksesibilitas,
			CreatedAt:           job.CreatedAt,
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

func (s *jobBoardService) GetByID(id uuid.UUID) (*dto.JobListingDetailResponse, error) {
	job, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, errors.New("job not found")
	}

	company, _ := s.repo.FindCompanyByID(job.CompanyID)
	companyName := ""
	companyLogo := ""
	if company != nil {
		companyName = company.Nama
		companyLogo = company.LogoURL
	}

	return &dto.JobListingDetailResponse{
		JobListingResponse: dto.JobListingResponse{
			ID:                  job.ID,
			CompanyID:           job.CompanyID,
			CompanyName:         companyName,
			CompanyLogo:         companyLogo,
			Judul:               job.Judul,
			Kota:                job.Kota,
			TipePekerjaan:       job.TipePekerjaan,
			BidangKerja:         job.BidangKerja,
			Gaji:                job.Gaji,
			DisabilitasDiterima: job.DisabilitasDiterima,
			LabelAksesibilitas:  job.LabelAksesibilitas,
			CreatedAt:           job.CreatedAt,
		},
		Deskripsi:   job.Deskripsi,
		Kualifikasi: job.Kualifikasi,
	}, nil
}

func (s *jobBoardService) ToggleSave(userID, jobID uuid.UUID) (bool, error) {
	// cek job exists
	job, err := s.repo.FindByID(jobID)
	if err != nil {
		return false, err
	}
	if job == nil {
		return false, errors.New("job not found")
	}

	// cek sudah disave atau belum
	isSaved, err := s.repo.IsJobSaved(userID, jobID)
	if err != nil {
		return false, err
	}

	if isSaved {
		// unsave
		if err := s.repo.UnsaveJob(userID, jobID); err != nil {
			return false, errors.New("failed to unsave job")
		}
		return false, nil
	}

	// save
	if err := s.repo.SaveJob(userID, jobID); err != nil {
		return false, errors.New("failed to save job")
	}
	return true, nil
}

func (s *jobBoardService) GetSavedJobs(userID uuid.UUID) ([]dto.JobListingResponse, error) {
	jobs, err := s.repo.FindSavedJobs(userID)
	if err != nil {
		return nil, errors.New("failed to get saved jobs")
	}

	var result []dto.JobListingResponse
	for _, job := range jobs {
		company, _ := s.repo.FindCompanyByID(job.CompanyID)
		companyName := ""
		companyLogo := ""
		if company != nil {
			companyName = company.Nama
			companyLogo = company.LogoURL
		}

		result = append(result, dto.JobListingResponse{
			ID:                  job.ID,
			CompanyID:           job.CompanyID,
			CompanyName:         companyName,
			CompanyLogo:         companyLogo,
			Judul:               job.Judul,
			Kota:                job.Kota,
			TipePekerjaan:       job.TipePekerjaan,
			BidangKerja:         job.BidangKerja,
			Gaji:                job.Gaji,
			DisabilitasDiterima: job.DisabilitasDiterima,
			LabelAksesibilitas:  job.LabelAksesibilitas,
			CreatedAt:           job.CreatedAt,
		})
	}

	return result, nil
}
