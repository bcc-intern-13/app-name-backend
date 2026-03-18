package entity

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	JobID           uuid.UUID `gorm:"type:uuid;not null;index" json:"job_id"`
	CvURL           string    `gorm:"column:cv_url;type:text;not null" json:"cv_url"`
	PortfolioLink   string    `gorm:"column:portfolio_link;type:text" json:"portfolio_link"`
	Status          string    `gorm:"column:status;type:varchar(20);not null;default:'Terkirim'" json:"status"`
	InterviewLink   string    `gorm:"column:interview_link;type:text" json:"interview_link"`
	InterviewPdfURL string    `gorm:"column:interview_pdf_url;type:text" json:"interview_pdf_url"`
	InterviewNotes  string    `gorm:"column:interview_notes;type:varchar(500)" json:"interview_notes"`
	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
