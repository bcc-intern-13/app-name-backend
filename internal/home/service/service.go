package service

import (
	"time"

	careerMappingDto "github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/home/dto"
	jobDto "github.com/bcc-intern-13/app-name-backend/internal/job_board/dto"
	onboardingDto "github.com/bcc-intern-13/app-name-backend/internal/onboarding/dto"
	"github.com/google/uuid"
)

type HomeService interface {
	GetSummary(userID uuid.UUID) (*dto.HomeSummaryResponse, error)
}

type homeService struct {
	onboardingRepo   onboardingDto.OnboardingRepository
	jobBoardService  jobDto.JobBoardService
	careerMappingSvc careerMappingDto.CareerMappingService
}

func NewHomeService(
	onboardingRepo onboardingDto.OnboardingRepository,
	jobBoardService jobDto.JobBoardService,
	careerMappingSvc careerMappingDto.CareerMappingService,
) HomeService {
	return &homeService{
		onboardingRepo:   onboardingRepo,
		jobBoardService:  jobBoardService,
		careerMappingSvc: careerMappingSvc,
	}
}

func (s *homeService) GetSummary(userID uuid.UUID) (*dto.HomeSummaryResponse, error) {
	profile, err := s.onboardingRepo.FindByUserID(userID)
	nama := ""
	if err == nil && profile != nil {
		nama = profile.Nama
	}

	greeting := dto.GreetingResponse{
		Nama:      nama,
		Timestamp: time.Now().UTC(),
	}

	var careerMapping *careerMappingDto.CareerMappingResponse
	cmResult, err := s.careerMappingSvc.GetLatestResult(userID)
	if err == nil {
		careerMapping = cmResult
	}

	var rekomendasi []jobDto.JobListingResponse
	if profile != nil {
		filter := jobDto.JobBoardFilter{
			BidangKerja:   profile.BidangKerja,
			TipePekerjaan: profile.TipePekerjaan,
			Limit:         5,
			Page:          1,
		}

		if careerMapping != nil && len(careerMapping.TopCategories) > 0 {
			filter.BidangKerja = mapCategoryToField(careerMapping.TopCategories[0].Code)
		}

		result, err := s.jobBoardService.GetAll(filter, userID)
		if err == nil && result != nil {
			rekomendasi = result.Data
		}
	}

	return &dto.HomeSummaryResponse{
		Greeting:            greeting,
		RekomendasiLowongan: rekomendasi,
		CareerMapping:       careerMapping,
	}, nil
}

func mapCategoryToField(code string) string {
	mapping := map[string]string{
		"KR": "Desain & Kreatif",
		"TK": "Teknologi & IT",
		"KO": "Administrasi & Keuangan",
		"ED": "Pendidikan",
		"AD": "Administrasi & Keuangan",
		"OP": "Administrasi & Keuangan",
	}
	if field, ok := mapping[code]; ok {
		return field
	}
	return ""
}
