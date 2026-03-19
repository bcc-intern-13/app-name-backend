package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/app/job_board/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/job_board/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/job_board/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type jobBoardRepository struct {
	db *gorm.DB
}

func NewJobBoardRepository(db *gorm.DB) contract.JobBoardRepository {
	return &jobBoardRepository{db: db}
}

func (r *jobBoardRepository) FindAll(filter dto.JobBoardFilter) ([]dto.JobListingWithCompany, int64, error) {
	var results []dto.JobListingWithCompany
	var total int64

	query := r.db.Table("job_listings jl").
		Select(`
			jl.*,
			c.name as company_name,
			c.logo_url as company_logo
		`).
		Joins("LEFT JOIN companies c ON c.id = jl.company_id").
		Where("jl.is_active = ?", true)

	if filter.City != "" {
		query = query.Where("jl.city ILIKE ?", "%"+filter.City+"%")
	}
	if filter.JobField != "" {
		query = query.Where("jl.job_field = ?", filter.JobField)
	}
	if filter.JobType != "" {
		query = query.Where("jl.job_type = ?", filter.JobType)
	}
	if filter.Disability != "" {
		query = query.Where("jl.accepted_disability @> ?", `["`+filter.Disability+`"]`)
	}
	if filter.AccessibilityLabel != "" {
		query = query.Where("jl.accessibility_label @> ?", `["`+filter.AccessibilityLabel+`"]`)
	}
	if filter.Search != "" {
		query = query.Where("jl.title ILIKE ?", "%"+filter.Search+"%")
	}

	query.Count(&total)

	page := filter.Page
	limit := filter.Limit
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := query.Order("jl.created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	return results, total, err
}
func (r *jobBoardRepository) FindByID(id uuid.UUID) (*entity.JobListing, error) {
	var job entity.JobListing
	err := r.db.Where("id = ?", id).First(&job).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &job, err
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

func (r *jobBoardRepository) FindSavedJobs(userID uuid.UUID) ([]dto.JobListingWithCompany, error) {
	var results []dto.JobListingWithCompany
	err := r.db.Table("job_listings jl").
		Select(`
			jl.*,
			c.name as company_name,
			c.logo_url as company_logo
		`).
		Joins("LEFT JOIN companies c ON c.id = jl.company_id").
		Joins("INNER JOIN saved_jobs sj ON sj.job_id = jl.id").
		Where("sj.user_id = ? AND jl.is_active = true", userID).
		Order("sj.created_at DESC").
		Scan(&results).Error
	return results, err
}
