package entity

import (
	"time"

	"github.com/google/uuid"
)

type SavedJob struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index"`
	JobID     uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
