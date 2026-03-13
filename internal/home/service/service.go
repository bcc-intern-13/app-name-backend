package service

import (
	"log/slog"
	"time"

	careerMappingDto "github.com/bcc-intern-13/app-name-backend/internal/career_mapping/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/home/dto"
	jobDto "github.com/bcc-intern-13/app-name-backend/internal/job_board/dto"
	onboardingDto "github.com/bcc-intern-13/app-name-backend/internal/onboarding/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type HomeService interface {
	GetSummary(userID uuid.UUID) (*dto.HomeSummaryResponse, *response.APIError)
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

func (s *homeService) GetSummary(userID uuid.UUID) (*dto.HomeSummaryResponse, *response.APIError) {
	//get profile, not fatal if error name just will be empyy
	profile, err := s.onboardingRepo.FindByUserID(userID)
	nama := ""
	if err != nil {
		slog.Error("failed to get user profile for home summary", "error", err, "userID", userID)
	} else if profile != nil {
		nama = profile.Nama
	}

	greeting := dto.GreetingResponse{
		Nama:      nama,
		Timestamp: time.Now().UTC(),
	}
	//get career mapping
	var careerMapping *careerMappingDto.CareerMappingResponse
	cmResult, apiErr := s.careerMappingSvc.GetLatestResult(userID)
	if apiErr != nil {
		//user havent done career mapping test
		if apiErr.Status != 404 {
			slog.Error("failed to get career mapping result", "error", apiErr.Message, "userID", userID)
		}
	} else {
		careerMapping = cmResult
	}

	// get job recommendations, not fatal if error, just return empty list
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

		result, apiErr := s.jobBoardService.GetAll(filter, userID)
		if apiErr != nil {
			slog.Error("failed to get job recommendations", "error", apiErr.Message, "userID", userID)
		} else if result != nil {
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
