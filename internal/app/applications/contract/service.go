package contract

import (
	"mime/multipart"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/applications/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type ApplicationService interface {
	Submit(userID uuid.UUID, req *dto.SubmitApplicationRequest, cv *multipart.FileHeader) *response.APIError
	GetAll(userID uuid.UUID, status string) ([]dto.ApplicationResponse, *response.APIError)
	GetByID(userID, id uuid.UUID) (*dto.ApplicationDetailResponse, *response.APIError)
	Delete(userID, id uuid.UUID) *response.APIError
}
