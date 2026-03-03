package database

import (
	domain "app-name/internal/domain/entity"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		log.Println("Failed to Auto Migrate to database.")
	}
	log.Println("Migration Completed....")
}
