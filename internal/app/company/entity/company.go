package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Company struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name               string         `gorm:"column:name;type:varchar(200);not null" json:"name"`
	LogoURL            string         `gorm:"column:logo_url;type:text" json:"logo_url"`
	Description        string         `gorm:"column:description;type:text" json:"description"`
	Industry           string         `gorm:"column:industry;type:varchar(100)" json:"industry"`
	Size               string         `gorm:"column:size;type:varchar(50)" json:"size"`
	Location           string         `gorm:"column:location;type:varchar(100)" json:"location"`
	Website            string         `gorm:"column:website;type:text" json:"website"`
	AcceptedDisability datatypes.JSON `gorm:"column:accepted_disability;type:jsonb;not null" json:"accepted_disability"`
	AccessibilityLabel datatypes.JSON `gorm:"column:accessibility_label;type:jsonb;not null" json:"accessibility_label"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
}
