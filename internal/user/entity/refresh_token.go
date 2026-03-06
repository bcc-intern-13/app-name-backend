package entity

import (
	"time"

	"github.com/google/uuid"
)

type RefreshToken struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `gorm:"type:text;not null;uniqueIndex"`
	ExpiredAt time.Time `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
