package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CV struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID         uuid.UUID      `gorm:"type:uuid;not null;uniqueIndex" json:"user_id"`
	Summary        string         `gorm:"column:summary;type:text" json:"summary"`
	Education      datatypes.JSON `gorm:"column:education;type:jsonb" json:"education"`
	Experience     datatypes.JSON `gorm:"column:experience;type:jsonb" json:"experience"`
	Skills         datatypes.JSON `gorm:"column:skills;type:jsonb" json:"skills"`
	AdaptiveSkills datatypes.JSON `gorm:"column:adaptive_skills;type:jsonb" json:"adaptive_skills"`
	CvScore        int            `gorm:"column:cv_score;default:0" json:"cv_score"`
	LastJobMatch   datatypes.JSON `gorm:"column:last_job_match;type:jsonb" json:"last_job_match"`
	IsAiVerified   bool           `gorm:"column:is_ai_verified;default:false" json:"is_ai_verified"`
	AiCallsToday   int            `gorm:"column:ai_calls_today;default:0" json:"ai_calls_today"`
	UpdatedAt      time.Time      `gorm:"autoUpdateTime" json:"updated_at"`

	CvURL string `gorm:"column:cv_url;type:text" json:"cv_url"`
}
