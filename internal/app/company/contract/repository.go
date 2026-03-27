package contract

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/company/entity"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
	"github.com/google/uuid"
)

type CompanyRepository interface {
	FindCompanyByID(id uuid.UUID) (*entity.Company, error)
	FindActiveJobsByCompanyID(id uuid.UUID) ([]dto.JobListingWithCompany, error)
}
