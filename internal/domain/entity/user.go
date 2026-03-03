package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Email      string    `gorm:"type:varchar(255);uniqueIndex;not null"         json:"email"`
	Password   string    `gorm:"type:text;not null"                             json:"-"`
	AvatarURL  string    `json:"avatar_url" gorm:"type:varchar(255);not null"`
	IsVerified bool      `json:"is_verified" gorm:"type:boolean;default:false;not null"`

	CreatedAt time.Time      `gorm:"autoCreateTime"                                 json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"                                 json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                          json:"-"`
}

//will be using this for future features.
// type Profile struct {
// 	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
// 	UserID    uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"                 json:"user_id"`
// 	FullName  string    `gorm:"type:varchar(100)"                              json:"full_name"`
// 	Phone     string    `gorm:"type:varchar(20)"                               json:"phone,omitempty"`
// 	Bio       string    `gorm:"type:text"                                      json:"bio,omitempty"`
// 	AvatarURL string    `gorm:"type:text"                                      json:"avatar_url,omitempty"`
// 	UpdatedAt time.Time `gorm:"autoUpdateTime"                                 json:"updated_at"`
// }
