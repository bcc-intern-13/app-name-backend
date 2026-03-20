package contract

import (
	"context"

	"github.com/bcc-intern-13/app-name-backend/internal/app/gemini/entity"
	"github.com/google/uuid"
)

type CVRepository interface {
	FindByUserID(ctx context.Context, userID uuid.UUID) (*entity.CV, error)
	Create(ctx context.Context, cv *entity.CV) error
	Update(ctx context.Context, cv *entity.CV) error
}
