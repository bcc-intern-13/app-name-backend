package contract

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/entity"
	"github.com/google/uuid"
)

type ApplicationRepository interface {
	Create(application *entity.Application) error
	FindAllByUserID(userID uuid.UUID, status string) ([]dto.ApplicationWithJob, error)
	FindByID(id uuid.UUID) (*dto.ApplicationWithJob, error)
	FindByUserIDAndJobID(userID, jobID uuid.UUID) (*entity.Application, error)
	FindLatestWithCVByUserID(userID uuid.UUID) (*entity.Application, error)
	Delete(id uuid.UUID) error
}
