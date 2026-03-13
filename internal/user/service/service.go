package service

import (
	"log/slog"
	"time"

	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/email"
	jwt "github.com/bcc-intern-13/app-name-backend/pkg/jwt"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userAuthService struct {
	repo                  dto.UserRepository
	refreshTokenRepo      dto.RefreshTokenRepository
	jwtSecret             string
	email                 *email.EmailService
	verificationTokenRepo dto.VerificationTokenRepository
}

func NewUserAuthService(
	repo dto.UserRepository,
	jwtSecret string,
	refreshTokenRepo dto.RefreshTokenRepository,
	verificationTokenRepo dto.VerificationTokenRepository,
	email *email.EmailService,
) dto.UserAuthService {
	return &userAuthService{
		repo:                  repo,
		refreshTokenRepo:      refreshTokenRepo,
		verificationTokenRepo: verificationTokenRepo,
		email:                 email,
		jwtSecret:             jwtSecret,
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

	refreshTokenStr := jwt.GenerateRefreshToken()
	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour), //one week
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
	//find refresh token
	refreshToken, err := s.refreshTokenRepo.FindByToken(token)
	if err != nil {
		slog.Error("failed to find refresh token", "error", err)
		return nil, response.ErrInternal("failed to find token")
	}
	if refreshToken == nil {
		return nil, response.ErrUnAuthorized("session expired, please login again")
	}
	//check if refresh token is expired
	if time.Now().After(refreshToken.ExpiredAt) {
		return nil, response.ErrUnAuthorized("session expired, please login again")
	}

	//find token by user id

	user, err := s.repo.FindByID(refreshToken.UserID.String())
	if err != nil {
		slog.Error("failed to find user by id", "error", err, "userID", refreshToken.UserID)
		return nil, response.ErrInternal("failed to find user")
	}
	if user == nil {
		return nil, response.ErrUnAuthorized("session expired, please login again")
	}

	//generate new access token
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

// vertification email
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
