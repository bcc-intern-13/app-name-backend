package dto

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationResponse struct {
	ID            uuid.UUID `json:"id"`
	JobID         uuid.UUID `json:"job_id"`
	JobTitle      string    `json:"job_title"`
	CompanyName   string    `json:"company_name"`
	CompanyLogo   string    `json:"company_logo"`
	Status        string    `json:"status"`
	PortfolioLink string    `json:"portfolio_link"`
	CvURL         string    `json:"cv_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type ApplicationDetailResponse struct {
	ApplicationResponse
	InterviewLink   string `json:"interview_link"`
	InterviewPdfURL string `json:"interview_pdf_url"`
	InterviewNotes  string `json:"interview_notes"`
}
