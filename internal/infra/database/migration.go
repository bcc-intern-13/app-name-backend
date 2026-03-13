package database

import (
	"log"

	careermappingentity "github.com/bcc-intern-13/app-name-backend/internal/career_mapping/entity"
	jobboardidentity "github.com/bcc-intern-13/app-name-backend/internal/job_board/entity"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&entity.User{},
		&entity.RefreshToken{},
		&entity.VerificationToken{},
		&entity.UserProfile{},
		&careermappingentity.CareerMappingQuestion{},
		&careermappingentity.CareerMappingResult{},
		&careermappingentity.CareerCategory{},
		&jobboardidentity.Company{},
		&jobboardidentity.JobListing{},
		&jobboardidentity.SavedJob{},
	)
	if err != nil {
		log.Println("Failed to Auto Migrate to database.")
	}
	log.Println("Migration Completed....")
}
