package repository

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/career_mapping/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type careerMappingRepository struct {
	db *gorm.DB
}

func NewCareerMappingRepository(db *gorm.DB) dto.CareerMappingRepository {
	return &careerMappingRepository{db: db}
}

func (r *careerMappingRepository) GetAllQuestions() ([]entity.CareerMappingQuestion, error) {
	var questions []entity.CareerMappingQuestion
	err := r.db.Order("nomor asc").Find(&questions).Error
	return questions, err
}

func (r *careerMappingRepository) CreateResult(result *entity.CareerMappingResult) error {
	return r.db.Create(result).Error
}

func (r *careerMappingRepository) FindLatestResultByUserID(userID uuid.UUID) (*entity.CareerMappingResult, error) {
	var result entity.CareerMappingResult
	err := r.db.Where("user_id = ?", userID).
		Order("attempt_number desc").
		First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, err
}

func (r *careerMappingRepository) CountByUserID(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&entity.CareerMappingResult{}).
		Where("user_id = ?", userID).
		Count(&count).Error
	return count, err
}

func (r *careerMappingRepository) GetCategoryByID(id string) (*entity.CareerCategory, error) {
	var category entity.CareerCategory
	err := r.db.Where("id = ?", id).First(&category).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}
