package repository

import (
	"errors"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/entity"
	jobDto "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type companyRepository struct {
	db *gorm.DB
}

func NewCompanyRepository(db *gorm.DB) contract.CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) FindCompanyByID(id uuid.UUID) (*entity.Company, error) {
	var company entity.Company
	err := r.db.Where("id = ?", id).First(&company).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &company, err
}

func (r *companyRepository) FindActiveJobsByCompanyID(id uuid.UUID) ([]jobDto.JobListingWithCompany, error) {
	var results []jobDto.JobListingWithCompany
	err := r.db.Table("job_listings jl").
		Select(`
			jl.*,
			c.name as company_name,
			c.logo_url as company_logo
		`).
		Joins("LEFT JOIN companies c ON c.id = jl.company_id").
		Where("jl.company_id = ? AND jl.is_active = true", id).
		Order("jl.created_at desc").
		Scan(&results).Error
	return results, err
}
