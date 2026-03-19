package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) contract.ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *entity.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) FindAllByUserID(userID uuid.UUID, status string) ([]dto.ApplicationWithJob, error) {
	var results []dto.ApplicationWithJob

	query := r.db.Table("applications a").
		Select(`
			a.*,
			jl.title as job_title,
			jl.city as job_city,
			jl.job_type as job_type,
			c.name as company_name,
			c.logo_url as company_logo
		`).
		Joins("LEFT JOIN job_listings jl ON jl.id = a.job_id").
		Joins("LEFT JOIN companies c ON c.id = jl.company_id").
		Where("a.user_id = ?", userID)

	if status != "" {
		query = query.Where("a.status = ?", status)
	}

	err := query.Order("a.created_at desc").Scan(&results).Error
	return results, err
}

func (r *applicationRepository) FindByID(id uuid.UUID) (*dto.ApplicationWithJob, error) {
	var result dto.ApplicationWithJob
	err := r.db.Table("applications a").
		Select(`
			a.*,
			jl.title as job_title,
			jl.city as job_city,
			jl.job_type as job_type,
			c.name as company_name,
			c.logo_url as company_logo
		`).
		Joins("LEFT JOIN job_listings jl ON jl.id = a.job_id").
		Joins("LEFT JOIN companies c ON c.id = jl.company_id").
		Where("a.id = ?", id).
		First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, err
}

func (r *applicationRepository) FindByUserIDAndJobID(userID, jobID uuid.UUID) (*entity.Application, error) {
	var application entity.Application
	err := r.db.Where("user_id = ? AND job_id = ?", userID, jobID).First(&application).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &application, err
}

func (r *applicationRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&entity.Application{}).Error
}
