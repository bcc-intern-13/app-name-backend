package database

import (
	"log"

	applicationsentity "github.com/bcc-intern-13/app-name-backend/internal/app/applications/entity"
	careermappingentity "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/entity"
	companyentity "github.com/bcc-intern-13/app-name-backend/internal/app/company/entity"
	cvEntity "github.com/bcc-intern-13/app-name-backend/internal/app/gemini/entity"
	jobboardidentity "github.com/bcc-intern-13/app-name-backend/internal/app/job_board/entity"
	orderentity "github.com/bcc-intern-13/app-name-backend/internal/app/payment/entity"
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/entity"

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
		&companyentity.Company{},
		&jobboardidentity.JobListing{},
		&jobboardidentity.SavedJob{},
		&applicationsentity.Application{},
		&cvEntity.CV{},
		&orderentity.Order{},
	)
	if err != nil {
		log.Println("Failed to Auto Migrate to database.")
	}
	log.Println("Migration Completed....")
}
