package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type CareerMappingQuestion struct {
	ID         uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Nomor      int            `gorm:"not null;uniqueIndex" json:"nomor"`
	Pertanyaan string         `gorm:"type:text;not null" json:"pertanyaan"`
	PilihanA   string         `gorm:"type:text;not null" json:"pilihan_a"`
	PilihanB   string         `gorm:"type:text;not null" json:"pilihan_b"`
	PilihanC   string         `gorm:"type:text;not null" json:"pilihan_c"`
	PilihanD   string         `gorm:"type:text;not null" json:"pilihan_d"`
	SkorA      datatypes.JSON `gorm:"type:jsonb;not null" json:"skor_a"`
	SkorB      datatypes.JSON `gorm:"type:jsonb;not null" json:"skor_b"`
	SkorC      datatypes.JSON `gorm:"type:jsonb;not null" json:"skor_c"`
	SkorD      datatypes.JSON `gorm:"type:jsonb;not null" json:"skor_d"`
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
	ID            string         `gorm:"type:varchar(2);primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(50);not null" json:"name"`
	Description   string         `gorm:"type:text;not null" json:"description"`
	JobsFormal    datatypes.JSON `gorm:"type:jsonb;not null" json:"jobs_formal"`
	JobsWirausaha datatypes.JSON `gorm:"type:jsonb;not null" json:"jobs_wirausaha"`
}
