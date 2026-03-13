package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/job_board/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/job_board/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type jobBoardRepository struct {
	db *gorm.DB
}

func NewJobBoardRepository(db *gorm.DB) dto.JobBoardRepository {
	return &jobBoardRepository{db: db}
}

func (r *jobBoardRepository) FindAll(filter dto.JobBoardFilter) ([]entity.JobListing, int64, error) {
	var jobs []entity.JobListing
	var total int64

	query := r.db.Model(&entity.JobListing{}).Where("is_active = ?", true)

	if filter.Kota != "" {
		query = query.Where("kota ILIKE ?", "%"+filter.Kota+"%")
	}
	if filter.BidangKerja != "" {
		query = query.Where("bidang_kerja = ?", filter.BidangKerja)
	}
	if filter.TipePekerjaan != "" {
		query = query.Where("tipe_pekerjaan = ?", filter.TipePekerjaan)
	}
	if filter.Disabilitas != "" {
		query = query.Where("disabilitas_diterima @> ?", `["`+filter.Disabilitas+`"]`)
	}
	if filter.LabelAksesibilitas != "" {
		query = query.Where("label_aksesibilitas @> ?", `["`+filter.LabelAksesibilitas+`"]`)
	}
	if filter.Search != "" {
		query = query.Where("judul ILIKE ?", "%"+filter.Search+"%")
	}

	query.Count(&total)

	// pagination
	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&jobs).Error

	return jobs, total, err
}

func (r *jobBoardRepository) FindByID(id uuid.UUID) (*entity.JobListing, error) {
	var job entity.JobListing
	err := r.db.Where("id = ?", id).First(&job).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &job, err
}

func (r *jobBoardRepository) FindCompanyByID(id uuid.UUID) (*entity.Company, error) {
	var company entity.Company
	err := r.db.Where("id = ?", id).First(&company).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &company, err
}

func (r *jobBoardRepository) SaveJob(userID, jobID uuid.UUID) error {
	return r.db.Exec(
		"INSERT INTO saved_jobs (id, user_id, job_id, created_at) VALUES (gen_random_uuid(), ?, ?, NOW())",
		userID, jobID,
	).Error
}

func (r *jobBoardRepository) UnsaveJob(userID, jobID uuid.UUID) error {
	return r.db.Exec(
		"DELETE FROM saved_jobs WHERE user_id = ? AND job_id = ?",
		userID, jobID,
	).Error
}

func (r *jobBoardRepository) IsJobSaved(userID, jobID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Raw(
		"SELECT COUNT(*) FROM saved_jobs WHERE user_id = ? AND job_id = ?",
		userID, jobID,
	).Scan(&count).Error
	return count > 0, err
}

func (r *jobBoardRepository) FindSavedJobs(userID uuid.UUID) ([]entity.JobListing, error) {
	var jobs []entity.JobListing
	err := r.db.Raw(`
		SELECT jl.* FROM job_listings jl
		INNER JOIN saved_jobs sj ON sj.job_id = jl.id
		WHERE sj.user_id = ? AND jl.is_active = true
		ORDER BY sj.created_at DESC
	`, userID).Scan(&jobs).Error
	return jobs, err
}
