package dto

import (
	"time"

	jobDto "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/dto"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CompanyResponse struct {
	ID                 uuid.UUID                   `json:"id"`
	Name               string                      `json:"name"`
	LogoURL            string                      `json:"logo_url"`
	Description        string                      `json:"description"`
	Industry           string                      `json:"industry"`
	Size               string                      `json:"size"`
	Location           string                      `json:"location"`
	Website            string                      `json:"website"`
	AcceptedDisability datatypes.JSON              `json:"accepted_disability"`
	AccessibilityLabel datatypes.JSON              `json:"accessibility_label"`
	CreatedAt          time.Time                   `json:"created_at"`
	JobListings        []jobDto.JobListingResponse `json:"job_listings"`
}
