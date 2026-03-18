package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/contract"
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

func (r *applicationRepository) FindAllByUserID(userID uuid.UUID, status string) ([]entity.Application, error) {
	var applications []entity.Application
	query := r.db.Where("user_id = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("created_at desc").Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) FindByID(id uuid.UUID) (*entity.Application, error) {
	var application entity.Application
	err := r.db.Where("id = ?", id).First(&application).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &application, err
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
