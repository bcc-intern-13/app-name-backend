package service

import (
	"log/slog"

	careerMappingContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/contract"
	careerMappingDto "github.com/bcc-intern-13/WorkAble-backend/internal/app/career_mapping/dto"
	onboardingContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/contract"
	userContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/smart_profile/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/smart_profile/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
)

type smartProfileService struct {
	onboardingRepo   onboardingContract.OnboardingRepository
	careerMappingSvc careerMappingContract.CareerMappingService
	userRepo         userContract.UserRepository
}

func NewSmartProfileService(
	onboardingRepo onboardingContract.OnboardingRepository,
	careerMappingSvc careerMappingContract.CareerMappingService,
	userRepo userContract.UserRepository,

) contract.SmartProfileService {
	return &smartProfileService{
		onboardingRepo:   onboardingRepo,
		careerMappingSvc: careerMappingSvc,
		userRepo:         userRepo,
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

	avatarURL := ""
	user, err := s.userRepo.FindByID(userID.String())
	if err != nil {
		slog.Error("failed to get user for avatar", "error", err, "userID", userID)
	} else if user != nil {
		avatarURL = user.AvatarURL
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
			AvatarURL: avatarURL,
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
