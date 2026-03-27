package service

import (
	"log/slog"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/onboarding/dto"
	userContract "github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type onboardingService struct {
	repo     contract.OnboardingRepository
	userRepo userContract.UserRepository
}

func NewOnboardingService(repo contract.OnboardingRepository, userRepo userContract.UserRepository) contract.OnboardingService {
	return &onboardingService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *onboardingService) Submit(userID uuid.UUID, req *dto.SubmitOnboardingRequest) *response.APIError {
	existing, err := s.repo.FindByUserID(userID)
	if err != nil {
		slog.Error("failed to check existing profile", "error", err, "userID", userID)
		return response.ErrInternal("failed to check existing profile")
	}
	if existing != nil {
		return response.ErrConflict("onboarding already completed")
	}

	//make new profile
	profile := &entity.UserProfile{
		ID:                      uuid.New(),
		UserID:                  userID,
		Name:                    req.Name,
		Age:                     req.Age,
		City:                    req.City,
		Education:               req.Education,
		JobField:                req.JobField,
		JobType:                 req.JobType,
		Status:                  req.Status,
		CommunicationPreference: req.CommunicationPreference,
		WorkEnvironment:         datatypes.JSON(req.WorkEnvironment),
		SpecialNeeds:            datatypes.JSON(req.SpecialNeeds),
	}

	if err := s.repo.Create(profile); err != nil {
		slog.Error("failed to save onboarding data", "error", err, "userID", userID)
		return response.ErrInternal("failed to save onboarding data")
	}

	if err := s.userRepo.UpdateOnboardingCompleted(userID); err != nil {
		slog.Error("failed to update onboarding status", "error", err, "userID", userID)
		return response.ErrInternal("failed to update onboarding status")
	}

	return nil
}

func (s *onboardingService) GetByUserID(userID uuid.UUID) (*entity.UserProfile, *response.APIError) {
	profile, err := s.repo.FindByUserID(userID)
	if err != nil {
		slog.Error("failed to get profile", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get profile")
	}
	if profile == nil {
		return nil, response.ErrNotFound("profile not found")
	}
	return profile, nil
}

func (s *onboardingService) Update(userID uuid.UUID, req *dto.SubmitOnboardingRequest) *response.APIError {
	profile, err := s.repo.FindByUserID(userID)
	if err != nil {
		slog.Error("failed to get profile", "error", err, "userID", userID)
		return response.ErrInternal("failed to get profile")
	}
	if profile == nil {
		return response.ErrNotFound("profile not found")
	}

	profile.Name = req.Name
	profile.Age = req.Age
	profile.City = req.City
	profile.Education = req.Education
	profile.JobField = req.JobField
	profile.JobType = req.JobType
	profile.Status = req.Status
	profile.CommunicationPreference = req.CommunicationPreference
	profile.WorkEnvironment = datatypes.JSON(req.WorkEnvironment)
	profile.SpecialNeeds = datatypes.JSON(req.SpecialNeeds)

	if err := s.repo.Update(profile); err != nil {
		slog.Error("failed to update profile", "error", err, "userID", userID)
		return response.ErrInternal("failed to update profile")
	}

	return nil
}
