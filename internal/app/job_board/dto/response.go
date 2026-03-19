package dto

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type JobListingResponse struct {
	ID                 uuid.UUID      `json:"id"`
	CompanyID          uuid.UUID      `json:"company_id"`
	CompanyName        string         `json:"company_name"`
	CompanyLogo        string         `json:"company_logo"`
	Title              string         `json:"title"`
	City               string         `json:"city"`
	JobType            string         `json:"job_type"`
	JobField           string         `json:"job_field"`
	Salary             string         `json:"salary"`
	AcceptedDisability datatypes.JSON `json:"accepted_disability"`
	AccessibilityLabel datatypes.JSON `json:"accessibility_label"`
	CreatedAt          time.Time      `json:"created_at"`
}

type JobListingDetailResponse struct {
	JobListingResponse
	Description   string `json:"description"`
	Qualification string `json:"qualification"`
}

type PaginatedJobResponse struct {
	Data  []JobListingResponse `json:"data"`
	Total int64                `json:"total"`
	Page  int                  `json:"page"`
	Limit int                  `json:"limit"`
}

type JobListingWithCompany struct {
	ID                 uuid.UUID      `json:"id"`
	CompanyID          uuid.UUID      `json:"company_id"`
	CompanyName        string         `json:"company_name"`
	CompanyLogo        string         `json:"company_logo"`
	Title              string         `json:"title"`
	City               string         `json:"city"`
	JobType            string         `json:"job_type"`
	JobField           string         `json:"job_field"`
	Salary             string         `json:"salary"`
	AcceptedDisability datatypes.JSON `json:"accepted_disability"`
	AccessibilityLabel datatypes.JSON `json:"accessibility_label"`
	IsActive           bool           `json:"is_active"`
	CreatedAt          time.Time      `json:"created_at"`
}
