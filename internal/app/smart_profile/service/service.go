package service

import (
	"log/slog"

	careerMappingContract "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/contract"
	careerMappingDto "github.com/bcc-intern-13/app-name-backend/internal/app/career_mapping/dto"
	onboardingContract "github.com/bcc-intern-13/app-name-backend/internal/app/onboarding/contract"

	"github.com/bcc-intern-13/app-name-backend/internal/app/smart_profile/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/smart_profile/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
)

type smartProfileService struct {
	onboardingRepo   onboardingContract.OnboardingRepository
	careerMappingSvc careerMappingContract.CareerMappingService
}

func NewSmartProfileService(
	onboardingRepo onboardingContract.OnboardingRepository,
	careerMappingSvc careerMappingContract.CareerMappingService,
) contract.SmartProfileService {
	return &smartProfileService{
		onboardingRepo:   onboardingRepo,
		careerMappingSvc: careerMappingSvc,
	}
}

func (s *smartProfileService) GetByUserID(userID uuid.UUID) (*dto.SmartProfileResponse, *response.APIError) {

	//get onboarding data by user id
	profile, err := s.onboardingRepo.FindByUserID(userID)
	if err != nil {
		slog.Error("failed to get user profile", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get profile")
	}
	if profile == nil {
		return nil, response.ErrNotFound("profile not found, please complete onboarding first")
	}

	// get newest career mapping result by user id
	var careerMapping *careerMappingDto.CareerMappingResponse
	cmResult, apiErr := s.careerMappingSvc.GetLatestResult(userID)
	if apiErr != nil {
		if apiErr.Status != 404 {
			slog.Error("failed to get career mapping", "error", apiErr.Message, "userID", userID)
		}
	} else {
		careerMapping = cmResult
	}

	return &dto.SmartProfileResponse{
		PersonalInfo: dto.PersonalInfoResponse{
			Name:      profile.Name,
			Age:       profile.Age,
			City:      profile.City,
			Education: profile.Education,
		},
		CareerPreference: dto.CareerPreferenceResponse{
			JobField: profile.JobField,
			JobType:  profile.JobType,
			Status:   profile.Status,
		},
		Communication: dto.CommunicationResponse{
			CommunicationPreference: profile.CommunicationPreference,
		},
		Accessibility: dto.AccessibilityResponse{
			WorkEnvironment: profile.WorkEnvironment,
			SpecialNeeds:    profile.SpecialNeeds,
		},
		CareerMapping: careerMapping,
		UpdatedAt:     profile.UpdatedAt,
	}, nil
}
