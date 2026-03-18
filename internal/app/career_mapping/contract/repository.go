package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/entity"
	"github.com/google/uuid"
)

type CareerMappingRepository interface {
	GetAllQuestions() ([]entity.CareerMappingQuestion, error)
	CreateResult(result *entity.CareerMappingResult) error
	FindLatestResultByUserID(userID uuid.UUID) (*entity.CareerMappingResult, error)
	CountByUserID(userID uuid.UUID) (int64, error)
	GetCategoryByID(id string) (*entity.CareerCategory, error)
}
