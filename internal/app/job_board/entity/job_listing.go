package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type JobListing struct {
	ID                 uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	CompanyID          uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	Title              string         `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Description        string         `gorm:"column:description;type:text;not null" json:"description"`
	Qualification      string         `gorm:"column:qualification;type:text;not null" json:"qualification"`
	City               string         `gorm:"column:city;type:varchar(100);not null" json:"city"`
	JobType            string         `gorm:"column:job_type;type:varchar(20);not null" json:"job_type"`
	JobField           string         `gorm:"column:job_field;type:varchar(50);not null" json:"job_field"`
	Salary             string         `gorm:"column:salary;type:varchar(100)" json:"salary"`
	AcceptedDisability datatypes.JSON `gorm:"column:accepted_disability;type:jsonb;not null" json:"accepted_disability"`
	AccessibilityLabel datatypes.JSON `gorm:"column:accessibility_label;type:jsonb;not null" json:"accessibility_label"`
	IsActive           bool           `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time      `gorm:"autoCreateTime" json:"created_at"`
}

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
