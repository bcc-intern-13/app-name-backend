package service

import (
	"log/slog"
	"time"

	careerMappingContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/contract"
	careerMappingDto "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/dto"

	userContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/home/dto"
	jobBoardContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/contract"
	jobDto "github.com/bcc-intern-13/WorkAble-backend/internal/app/job_board/dto"
	onboardingContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/contract"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type HomeService interface {
	GetSummary(userID uuid.UUID) (*dto.HomeSummaryResponse, *response.APIError)
}

type homeService struct {
	onboardingRepo   onboardingContract.OnboardingRepository
	jobBoardService  jobBoardContract.JobBoardService
	careerMappingSvc careerMappingContract.CareerMappingService
	userRepo         userContract.UserRepository
}

func NewHomeService(
	onboardingRepo onboardingContract.OnboardingRepository,
	jobBoardService jobBoardContract.JobBoardService,
	careerMappingSvc careerMappingContract.CareerMappingService,
	userRepo userContract.UserRepository,
) HomeService {
	return &homeService{
		onboardingRepo:   onboardingRepo,
		jobBoardService:  jobBoardService,
		careerMappingSvc: careerMappingSvc,
		userRepo:         userRepo,
	}
}

func (s *homeService) GetSummary(userID uuid.UUID) (*dto.HomeSummaryResponse, *response.APIError) {
	//Get profile
	profile, err := s.onboardingRepo.FindByUserID(userID)
	nama := ""
	if err != nil {
		slog.Error("failed to get user profile for home summary", "error", err, "userID", userID)
	} else if profile != nil {
		nama = profile.Name
	}

	// Get Avatar
	avatarURL := ""
	user, err := s.userRepo.FindByID(userID.String())
	if err != nil {
		slog.Error("failed to get user for avatar", "error", err, "userID", userID)
	} else if user != nil {
		avatarURL = user.AvatarURL
	}

	greeting := dto.GreetingResponse{
		Name:      nama,
		Timestamp: time.Now().UTC(),
		AvatarURL: avatarURL,
	}

	var careerMapping *careerMappingDto.CareerMappingResponse
	cmResult, apiErr := s.careerMappingSvc.GetLatestResult(userID)
	if apiErr != nil {
		if apiErr.Status != 404 {
			slog.Error("failed to get career mapping result", "error", apiErr.Message, "userID", userID)
		}
	} else {
		careerMapping = cmResult
	}

	filter := jobDto.JobBoardFilter{
		Limit: 5,
		Page:  1,
	}

	if careerMapping != nil && len(careerMapping.TopCategories) > 0 {
		filter.JobField = mapCategoryToField(careerMapping.TopCategories[0].Code)

	} else if profile != nil {
		filter.JobField = profile.JobField
	}

	var jobRecommendations []jobDto.JobListingResponse
	result, apiErr := s.jobBoardService.GetAll(filter, userID)
	if apiErr != nil {
		slog.Error("failed to get job recommendations", "error", apiErr.Message, "userID", userID)
	} else if result != nil {
		jobRecommendations = result.Data
	}

	return &dto.HomeSummaryResponse{
		Greeting:           greeting,
		JobRecommendations: jobRecommendations,
		CareerMapping:      careerMapping,
	}, nil
}

func mapCategoryToField(code string) string {
	mapping := map[string]string{
		"KR": "Kreatif dan Seni",
		"TK": "Teknologi dan Digital",
		"KO": "Komunikasi dan Orang",
		"ED": "Edukasi dan Sosial",
		"AD": "Administrasi dan Data",
		"OP": "Operasional dan Detail",
	}

	if field, ok := mapping[code]; ok {
		return field
	}
	return ""
}
