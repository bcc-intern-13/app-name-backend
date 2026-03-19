package contract

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/company/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type CompanyService interface {
	GetByID(id uuid.UUID) (*dto.CompanyResponse, *response.APIError)
}
