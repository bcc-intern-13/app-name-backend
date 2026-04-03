package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name             string     `gorm:"type:varchar(100)" json:"nama"`
	Email            string     `gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"`
	Password         string     `gorm:"type:text;not null"                             json:"-"`
	AvatarURL        string     `json:"avatar_url" gorm:"type:varchar(255)"`
	IsVerified       bool       `json:"is_verified" gorm:"type:boolean;default:false;not null"`
	IsPremium        bool       `gorm:"default:false" json:"is_premium"`
	PremiumExpiresAt *time.Time `gorm:"column:premium_expires_at" json:"premium_expires_at"`

	OnboardingCompleted bool `gorm:"default:false" json:"onboarding_completed"`

	CreatedAt time.Time      `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                          json:"-"`

	ResetToken   *string    `gorm:"column:reset_token;type:varchar(255)" json:"-"`
	ResetExpires *time.Time `gorm:"column:reset_expires;type:timestamp" json:"-"`
}

type UserProfile struct {
	ID                      uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID                  uuid.UUID      `gorm:"type:uuid;not null" json:"user_id"`
	Name                    string         `gorm:"column:name;type:varchar(100);not null" json:"name"`
	Age                     string         `gorm:"column:age;type:varchar(20);not null" json:"age"`
	City                    string         `gorm:"column:city;type:varchar(100);not null" json:"city"`
	Education               string         `gorm:"column:education;type:varchar(20);not null" json:"education"`
	JobField                string         `gorm:"column:job_field;type:varchar(50);not null" json:"job_field"`
	JobType                 string         `gorm:"column:job_type;type:varchar(20);not null" json:"job_type"`
	Status                  string         `gorm:"column:status;type:varchar(30);not null" json:"status"`
	CommunicationPreference string         `gorm:"column:communication_preference;type:varchar(30);not null" json:"communication_preference"`
	WorkEnvironment         datatypes.JSON `gorm:"column:work_environment;type:jsonb" json:"work_environment"`
	SpecialNeeds            datatypes.JSON `gorm:"column:special_needs;type:jsonb" json:"special_needs"`
	UpdatedAt               time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	User                    User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
