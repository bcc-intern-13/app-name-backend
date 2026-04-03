package dto

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationWithJob struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	JobID           uuid.UUID  `json:"job_id"`
	JobTitle        string     `json:"job_title"`
	JobCity         string     `json:"job_city"`
	JobType         string     `json:"job_type"`
	CompanyName     string     `json:"company_name"`
	CompanyLogo     string     `json:"company_logo"`
	CvURL           string     `json:"cv_url"`
	PortfolioLink   string     `json:"portfolio_link"`
	Status          string     `json:"status"`
	InterviewLink   string     `json:"interview_link"`
	InterviewPdfURL string     `json:"interview_pdf_url"`
	InterviewNotes  string     `json:"interview_notes"`
	InterviewDate   *time.Time `json:"interviewDate"`
	Interviewer     string     `json:"interviewer"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
type ApplicationResponse struct {
	ID            uuid.UUID `json:"id"`
	JobID         uuid.UUID `json:"job_id"`
	JobTitle      string    `json:"job_title"`
	JobCity       string    `json:"job_city"`
	JobType       string    `json:"job_type"`
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
	InterviewLink   string     `json:"interview_link"`
	InterviewPdfURL string     `json:"interview_pdf_url"`
	InterviewNotes  string     `json:"interview_notes"`
	InterviewDate   *time.Time `json:"interviewDate"`
	Interviewer     string     `json:"interviewer"`
}
