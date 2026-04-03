package service

import (
	"context"
	"io"
	"log/slog"
	"mime/multipart"
	"strings"
	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/email"
	jwt "github.com/bcc-intern-13/WorkAble-backend/pkg/jwt"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/storage"
	util "github.com/bcc-intern-13/WorkAble-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type userAuthService struct {
	repo                  contract.UserRepository
	refreshTokenRepo      contract.RefreshTokenRepository
	jwtSecret             string
	email                 *email.EmailService
	verificationTokenRepo contract.VerificationTokenRepository
	storage               *storage.StorageService
	redisClient           *redis.Client
}

func NewUserAuthService(
	repo contract.UserRepository,
	jwtSecret string,
	refreshTokenRepo contract.RefreshTokenRepository,
	verificationTokenRepo contract.VerificationTokenRepository,
	email *email.EmailService,
	storageSvc *storage.StorageService,
	redisClient *redis.Client,
) contract.UserAuthService {
	return &userAuthService{
		repo:                  repo,
		refreshTokenRepo:      refreshTokenRepo,
		verificationTokenRepo: verificationTokenRepo,
		email:                 email,
		jwtSecret:             jwtSecret,
		storage:               storageSvc,
		redisClient:           redisClient,
	}
}

func (s *userAuthService) Register(req *dto.RegisterRequest) (*entity.User, *response.APIError) {
	if err := util.ValidatePassword(req.Password); err != nil {
		return nil, response.ErrBadRequest(err.Error())
	}

	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, response.ErrConflict("user is already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, response.ErrInternal("could not hash password")
	}

	user := &entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
	}

	verificationTokenStr := jwt.GenerateRefreshToken()
	verificationToken := &entity.VerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     verificationTokenStr,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.repo.RegisterTransaction(user, verificationToken); err != nil {
		slog.Error("failed to process registration transaction", "error", err)
		return nil, response.ErrInternal("failed to save registration data")
	}

	go s.email.SendVerificationEmail(user.Email, verificationTokenStr)

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
		ExpiredAt: time.Now().Add(refreshTokenDuration),
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		slog.Error("failed to save refresh token", "error", err, "userID", user.ID)
		return nil, response.ErrInternal("failed to save refresh token")
	}

	return &dto.LoginResponse{
		AccessToken:           accessToken,
		RefreshToken:          refreshTokenStr,
		RefreshTokenExpiresAt: refreshToken.ExpiredAt,
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
		AccessToken:           accessToken,
		RefreshToken:          token,
		RefreshTokenExpiresAt: refreshToken.ExpiredAt,
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

	//rate limiter
	ctx := context.Background()

	limitKey := "limit:resend_email:" + email

	isAllowed, err := s.redisClient.SetNX(ctx, limitKey, "1", 1*time.Minute).Result()
	if err != nil {
		slog.Error("Redis error on rate limit", "error", err)
	}

	if !isAllowed {
		return response.ErrBadRequest("Too many attemps, please wait 1 minute.")
	}

	//find user by gmail
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		slog.Error("failed to find user by email", "error", err, "email", email)
		return response.ErrInternal("failed to process request")
	}

	// user not found
	if user == nil {
		return response.ErrNotFound("user not found with this email")
	}

	// check if user isverified
	if user.IsVerified {
		return response.ErrBadRequest("user is already verified")
	}

	// generate new token
	verificationTokenStr := jwt.GenerateRefreshToken()
	verificationToken := &entity.VerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     verificationTokenStr,
		ExpiredAt: time.Now().Add(24 * time.Hour), // Berlaku 24 jam lagi
	}

	// save token to database
	if err := s.verificationTokenRepo.Create(verificationToken); err != nil {
		slog.Error("failed to save new verification token", "error", err, "userID", user.ID)
		return response.ErrInternal("failed to generate new verification link")
	}

	// resend the gmail
	if err := s.email.SendVerificationEmail(user.Email, verificationTokenStr); err != nil {
		slog.Error("failed to resend verification email", "error", err, "email", user.Email)
		return response.ErrInternal("failed to resend verification email")
	}

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

		user = &entity.User{
			ID:         uuid.New(),
			Email:      req.Email,
			Name:       req.Name,
			AvatarURL:  req.Picture,
			Password:   uuid.New().String(),
			IsVerified: true,
		}
		if err := s.repo.Create(user); err != nil {
			slog.Error("failed to create google user", "error", err)
			return nil, response.ErrInternal("failed to create user")
		}
	} else {
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
		AccessToken:           accessToken,
		RefreshToken:          refreshTokenStr,
		RefreshTokenExpiresAt: refreshToken.ExpiredAt,

		User: dto.UserData{
			ID:    user.ID.String(),
			Email: user.Email,
		},
	}, nil
}

func (s *userAuthService) UploadAvatar(ctx context.Context, userID uuid.UUID, file *multipart.FileHeader) (*dto.AvatarUploadResponse, *response.APIError) {
	// Validate formats
	contentType := file.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		return nil, response.ErrBadRequest("avatar must be an image (jpeg, png, webp)")
	}

	// Validate maximum capacity 2MB
	if file.Size > 2*1024*1024 {
		return nil, response.ErrBadRequest("avatar size must be less than 2MB")
	}

	// cek user old avvatar
	user, err := s.repo.FindByID(userID.String())
	if err != nil || user == nil {
		return nil, response.ErrInternal("failed to get user")
	}

	// erase old avatar
	if user.AvatarURL != "" {
		// cut avatar to gert relative path
		parts := strings.Split(user.AvatarURL, s.storage.BucketAvatar+"/")
		if len(parts) == 2 {
			oldFilePath := parts[1]
			_ = s.storage.DeleteFile(s.storage.BucketAvatar, oldFilePath)
			slog.Info("deleted old avatar from storage", "path", oldFilePath)
		}
	}

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

	avatarURL, err := s.storage.UploadAvatar(userID.String(), fileBytes, contentType)
	if err != nil {
		slog.Error("failed to upload avatar to storage", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to upload avatar")
	}

	user.AvatarURL = avatarURL
	if err := s.repo.Update(user); err != nil {
		slog.Error("failed to update user avatar in db", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to save avatar url")
	}

	return &dto.AvatarUploadResponse{AvatarURL: avatarURL}, nil
}

func (s *userAuthService) ForgotPassword(email string) *response.APIError {
	user, err := s.repo.FindByEmail(email)

	//keep sending succes to manipulate hacker
	if err != nil || user == nil {
		return nil
	}

	//unique token
	token := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	// store token to database
	user.ResetToken = &token
	user.ResetExpires = &expiresAt
	if err := s.repo.Update(user); err != nil {
		slog.Error("failed to save reset token", "error", err)
		return response.ErrInternal("Gagal memproses permintaan")
	}

	// send using email service
	if err := s.email.SendResetPasswordEmail(user.Email, token); err != nil {
		slog.Error("failed to send reset password email", "error", err, "email", user.Email)
		return response.ErrInternal("Gagal mengirim email reset password")
	}

	return nil
}

func (s *userAuthService) ResetPassword(token string, newPassword string) *response.APIError {

	user, err := s.repo.FindByResetToken(token)
	if err != nil || user == nil {
		return response.ErrBadRequest("Token tidak valid atau salah")
	}

	// check token expiry
	if user.ResetExpires != nil && time.Now().After(*user.ResetExpires) {
		return response.ErrBadRequest("Token sudah kadaluarsa, silakan request ulang")
	}

	if err := util.ValidatePassword(newPassword); err != nil {
		return response.ErrBadRequest(err.Error())
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("failed to hash new password", "error", err)
		return response.ErrInternal("Failed to create new password")
	}

	// Update password
	user.Password = string(hashedPassword)
	user.ResetToken = nil
	user.ResetExpires = nil

	if err := s.repo.Update(user); err != nil {
		slog.Error("failed to update user password", "error", err)
		return response.ErrInternal("Failed to reset user password")
	}

	return nil
}
