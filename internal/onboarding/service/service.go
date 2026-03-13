package service

import (
	"log/slog"

	"github.com/bcc-intern-13/app-name-backend/internal/onboarding/dto"
	userdto "github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type onboardingService struct {
	repo     dto.OnboardingRepository
	userRepo userdto.UserRepository
}

func NewOnboardingService(repo dto.OnboardingRepository, userRepo userdto.UserRepository) dto.OnboardingService {
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
		ID:                   uuid.New(),
		UserID:               userID,
		Nama:                 req.Nama,
		Usia:                 req.Usia,
		Kota:                 req.Kota,
		Pendidikan:           req.Pendidikan,
		BidangKerja:          req.BidangKerja,
		TipePekerjaan:        req.TipePekerjaan,
		Status:               req.Status,
		PreferensiKomunikasi: req.PreferensiKomunikasi,
		LingkunganKerja:      datatypes.JSON(req.LingkunganKerja),
		KebutuhanKhusus:      datatypes.JSON(req.KebutuhanKhusus),
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

	profile.Nama = req.Nama
	profile.Usia = req.Usia
	profile.Kota = req.Kota
	profile.Pendidikan = req.Pendidikan
	profile.BidangKerja = req.BidangKerja
	profile.TipePekerjaan = req.TipePekerjaan
	profile.Status = req.Status
	profile.PreferensiKomunikasi = req.PreferensiKomunikasi
	profile.LingkunganKerja = datatypes.JSON(req.LingkunganKerja)
	profile.KebutuhanKhusus = datatypes.JSON(req.KebutuhanKhusus)

	if err := s.repo.Update(profile); err != nil {
		slog.Error("failed to update profile", "error", err, "userID", userID)
		return response.ErrInternal("failed to update profile")
	}

	return nil
}
