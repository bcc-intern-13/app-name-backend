package service

import (
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/email"
	jwt "github.com/bcc-intern-13/WorkAble-backend/pkg/jwt"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/storage"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userAuthService struct {
	repo                  contract.UserRepository
	refreshTokenRepo      contract.RefreshTokenRepository
	jwtSecret             string
	email                 *email.EmailService
	verificationTokenRepo contract.VerificationTokenRepository
	storage               *storage.StorageService
}

func NewUserAuthService(
	repo contract.UserRepository,
	jwtSecret string,
	refreshTokenRepo contract.RefreshTokenRepository,
	verificationTokenRepo contract.VerificationTokenRepository,
	email *email.EmailService,
	storageSvc *storage.StorageService,
) contract.UserAuthService {
	return &userAuthService{
		repo:                  repo,
		refreshTokenRepo:      refreshTokenRepo,
		verificationTokenRepo: verificationTokenRepo,
		email:                 email,
		jwtSecret:             jwtSecret,
		storage:               storageSvc,
	}
}

func (s *userAuthService) Register(req *dto.RegisterRequest) (*entity.User, *response.APIError) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		slog.Error("failed to check existing user", "error", err)
		return nil, response.ErrInternal("failed to check existing user")
	}
	if existing != nil {
		return nil, response.ErrConflict("user is already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash password", "error", err)
		return nil, response.ErrInternal("could not hash password")
	}

	user := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := s.repo.Create(user); err != nil {
		slog.Error("failed to create user", "error", err, "email", req.Email)
		return nil, response.ErrInternal("failed to create user")
	}

	verificationTokenStr := jwt.GenerateRefreshToken()
	verificationToken := &entity.VerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     verificationTokenStr,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.verificationTokenRepo.Create(verificationToken); err != nil {
		slog.Error("failed to save verification token", "error", err, "userID", user.ID)
		return nil, response.ErrInternal("failed to save verification token")
	}

	if err := s.email.SendVerificationEmail(user.Email, verificationTokenStr); err != nil {
		slog.Error("failed to send verification email", "error", err, "email", user.Email)
		return nil, response.ErrInternal("failed to send verification email")
	}

	return user, nil
}

func (s *userAuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, *response.APIError) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		slog.Error("failed to find user by email", "error", err, "email", req.Email)
		return nil, response.ErrInternal("failed to find user")
	}
	if user == nil {
		return nil, response.ErrNotFound("user not found")
	}
	if !user.IsVerified {
		return nil, response.ErrUnAuthorized("please verify your email first")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, response.ErrUnAuthorized("Invalid Email or Password")
	}

	accessToken, err := jwt.GenerateToken(user, s.jwtSecret)
	if err != nil {
		slog.Error("failed to generate access token", "error", err, "userID", user.ID)
		return nil, response.ErrInternal("failed to generate token")
	}

	if err := s.refreshTokenRepo.DeleteByUserID(user.ID); err != nil {
		slog.Error("failed to delete old refresh tokens", "error", err)
		return nil, response.ErrInternal("failed to delete old refresh tokens")
	}

	//Generate refresh token
	//if rememmber me is on one week refresh time.
	//if off 24 hours refresh token only

	refreshTokenDuration := 24 * time.Hour

	if req.RememberMe {
		refreshTokenDuration = 7 * 24 * time.Hour
	}

	refreshTokenStr := jwt.GenerateRefreshToken()
	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiredAt: time.Now().Add(refreshTokenDuration), //one week
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		slog.Error("failed to save refresh token", "error", err, "userID", user.ID)
		return nil, response.ErrInternal("failed to save refresh token")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User: dto.UserData{
			ID:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}

func (s *userAuthService) RefreshToken(token string) (*dto.LoginResponse, *response.APIError) {
	//Find refresh token
	refreshToken, err := s.refreshTokenRepo.FindByToken(token)
	if err != nil {
		slog.Error("failed to find refresh token", "error", err)
		return nil, response.ErrInternal("failed to find token")
	}
	if refreshToken == nil {
		return nil, response.ErrUnAuthorized("session expired, please login again")
	}
	//Check if refresh token is expired
	if time.Now().After(refreshToken.ExpiredAt) {
		return nil, response.ErrUnAuthorized("session expired, please login again")
	}

	//Find token by user id
	user, err := s.repo.FindByID(refreshToken.UserID.String())
	if err != nil {
		slog.Error("failed to find user by id", "error", err, "userID", refreshToken.UserID)
		return nil, response.ErrInternal("failed to find user")
	}
	if user == nil {
		return nil, response.ErrUnAuthorized("session expired, please login again")
	}

	//Generate new access token
	accessToken, err := jwt.GenerateToken(user, s.jwtSecret)
	if err != nil {
		slog.Error("failed to generate access token", "error", err, "userID", user.ID)
		return nil, response.ErrInternal("failed to generate token")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: token,
		User: dto.UserData{
			ID:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}

func (s *userAuthService) Logout(token string) *response.APIError {
	if err := s.refreshTokenRepo.DeleteByToken(token); err != nil {
		slog.Error("failed to delete refresh token", "error", err)
		return response.ErrInternal("failed to delete refresh token")
	}
	return nil
}

// VerifyEmail verifies the user's email address
func (s *userAuthService) VerifyEmail(token string) *response.APIError {
	verificationToken, err := s.verificationTokenRepo.FindByToken(token)
	if err != nil {
		slog.Error("failed to find verification token", "error", err)
		return response.ErrInternal("failed to find verification token")
	}
	if verificationToken == nil {
		return response.ErrBadRequest("invalid or expired verification link")
	}
	if time.Now().After(verificationToken.ExpiredAt) {
		return response.ErrBadRequest("invalid or expired verification link")
	}

	//note updatge boolean user verified is using the dto interface of user repository not userautheservice.
	if err := s.repo.UpdateVerified(verificationToken.UserID); err != nil {
		slog.Error("failed to update verification status", "error", err, "userID", verificationToken.UserID)
		return response.ErrInternal("failed to update verification status")
	}

	if err := s.verificationTokenRepo.DeleteByToken(token); err != nil {
		slog.Error("failed to delete verification token", "error", err)
		return response.ErrInternal("failed to delete verification token")
	}

	return nil
}

func (s *userAuthService) ResendVerificationEmail(email string) *response.APIError {
	return nil
}

func (s *userAuthService) GoogleAuth(req *dto.GoogleAuthRequest) (*dto.LoginResponse, *response.APIError) {
	// cek apakah user sudah ada
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		slog.Error("failed to check existing user", "error", err, "email", req.Email)
		return nil, response.ErrInternal("failed to check existing user")
	}

	var user *entity.User

	if existing == nil {
		// register otomatis — tidak perlu password karena Google yang auth
		user = &entity.User{
			ID:         uuid.New(),
			Email:      req.Email,
			Nama:       req.Name,
			AvatarURL:  req.Picture,
			Password:   uuid.New().String(), // random password — tidak dipakai
			IsVerified: true,                // langsung verified karena dari Google
		}
		if err := s.repo.Create(user); err != nil {
			slog.Error("failed to create google user", "error", err)
			return nil, response.ErrInternal("failed to create user")
		}
	} else {
		// update avatar kalau sudah ada
		user = existing
		if req.Picture != "" && user.AvatarURL != req.Picture {
			user.AvatarURL = req.Picture
			s.repo.Update(user)
		}
	}

	// generate JWT
	accessToken, err := jwt.GenerateToken(user, s.jwtSecret)
	if err != nil {
		slog.Error("failed to generate token", "error", err)
		return nil, response.ErrInternal("failed to generate token")
	}

	// generate refresh token
	refreshTokenStr := jwt.GenerateRefreshToken()
	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiredAt: time.Now().Add(30 * 24 * time.Hour),
	}
	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		slog.Error("failed to save refresh token", "error", err)
		return nil, response.ErrInternal("failed to save refresh token")
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		User: dto.UserData{
			ID:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}

func (s *userAuthService) UploadAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.AvatarUploadResponse, *response.APIError) {
	// 1. Validasi Format (Cuma boleh gambar)
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return nil, response.ErrBadRequest("avatar must be an image (jpeg, png, webp)")
	}

	// 2. Validasi Ukuran (Maks 2MB)
	if file.Size > 2*1024*1024 {
		return nil, response.ErrBadRequest("avatar size must be less than 2MB")
	}

	// 3. Baca File
	f, err := file.Open()
	if err != nil {
		slog.Error("failed to open avatar file", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to process avatar file")
	}
	defer f.Close()

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		slog.Error("failed to read avatar file", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to read avatar file")
	}

	// 4. Lempar ke Kurir Supabase
	avatarURL, err := s.storage.UploadAvatar(userID.String(), fileBytes, contentType)
	if err != nil {
		slog.Error("failed to upload avatar to storage", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to upload avatar")
	}

	// 5. Update kolom avatar_url di tabel users
	user, err := s.repo.FindByID(userID.String())
	if err != nil || user == nil {
		return nil, response.ErrInternal("failed to get user")
	}

	user.AvatarURL = avatarURL
	// Asumsi gua lu punya fungsi Update di UserRepository
	if err := s.repo.Update(user); err != nil {
		slog.Error("failed to update user avatar", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to save avatar url")
	}

	return &dto.AvatarUploadResponse{AvatarURL: avatarURL}, nil
}
