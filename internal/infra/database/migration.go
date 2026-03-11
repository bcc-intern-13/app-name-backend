package database

import (
	"log"

	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.VerificationToken{},
		&entity.UserProfile{},
	)
	if err != nil {
		log.Println("Failed to Auto Migrate to database.")
	}
	log.Println("Migration Completed....")
}
