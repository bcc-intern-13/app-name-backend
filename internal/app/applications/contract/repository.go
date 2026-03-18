package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/applications/entity"
	"github.com/google/uuid"
)

type ApplicationRepository interface {
	Create(application *entity.Application) error
	FindAllByUserID(userID uuid.UUID, status string) ([]entity.Application, error)
	FindByID(id uuid.UUID) (*entity.Application, error)
	FindByUserIDAndJobID(userID, jobID uuid.UUID) (*entity.Application, error)
	Delete(id uuid.UUID) error
}
