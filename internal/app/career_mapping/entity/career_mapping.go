package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CareerMappingQuestion struct {
	ID       uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Number   int            `gorm:"not null;uniqueIndex" json:"number"`
	Question string         `gorm:"type:text;not null" json:"question"`
	OptionA  string         `gorm:"type:text;not null" json:"option_a"`
	OptionB  string         `gorm:"type:text;not null" json:"option_b"`
	OptionC  string         `gorm:"type:text;not null" json:"option_c"`
	OptionD  string         `gorm:"type:text;not null" json:"option_d"`
	ScoreA   datatypes.JSON `gorm:"type:jsonb;not null" json:"score_a"`
	ScoreB   datatypes.JSON `gorm:"type:jsonb;not null" json:"score_b"`
	ScoreC   datatypes.JSON `gorm:"type:jsonb;not null" json:"score_c"`
	ScoreD   datatypes.JSON `gorm:"type:jsonb;not null" json:"score_d"`
}

type CareerMappingResult struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Answers       datatypes.JSON `gorm:"type:jsonb;not null" json:"answers"`
	Scores        datatypes.JSON `gorm:"type:jsonb;not null" json:"scores"`
	TopCategories datatypes.JSON `gorm:"type:jsonb;not null" json:"top_categories"`
	AttemptNumber int            `gorm:"not null" json:"attempt_number"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

type CareerCategory struct {
	ID          string         `gorm:"type:varchar(2);primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(50);not null" json:"name"`
	Description string         `gorm:"type:text;not null" json:"description"`
	FormalJobs  datatypes.JSON `gorm:"type:jsonb;not null" json:"formal_jobs"`
	SideJobs    datatypes.JSON `gorm:"type:jsonb;not null" json:"side_jobs"`
}
