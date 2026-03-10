package service

import (
	"errors"
	"time"

	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/email"
	jwt "github.com/bcc-intern-13/app-name-backend/pkg/jwt"
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

func (s *userAuthService) Register(req *dto.RegisterRequest) (*entity.User, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("user is already registered")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		ID:       uuid.New(),
		Email:    req.Email,
		Password: string(hashed),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	//generate email verification link
	verificationTokenStr := jwt.GenerateRefreshToken() //we use same generation to gnenerate token
	verificationToken := &entity.VerificationToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     verificationTokenStr,
		ExpiredAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.verificationTokenRepo.Create(verificationToken); err != nil {
		return nil, errors.New("failed to save verification token")
	}

	if err := s.email.SendVerificationEmail(user.Email, verificationTokenStr); err != nil {
		return nil, errors.New("failed to send verification email")
	}

	return user, nil
}

func (s *userAuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {

	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	if !user.IsVerified {
		return nil, errors.New("email not verified yet.")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errors.New("Wrong password, please try again")
	}
	accessToken, err := jwt.GenerateToken(user, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	refreshTokenStr := jwt.GenerateRefreshToken()

	refreshToken := &entity.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     refreshTokenStr,
		ExpiredAt: time.Now().Add(7 * 24 * time.Hour), //one week
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, errors.New("failed to save refresh token")
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

func (s *userAuthService) RefreshToken(token string) (*dto.LoginResponse, error) {
	//find refresh token
	RefreshToken, err := s.refreshTokenRepo.FindByToken(token) // note ini fungsi s. buat apa ?

	if err != nil {
		return nil, errors.New("failed to find token")
	}
	if RefreshToken == nil {
		return nil, errors.New("refresh token not found") //todo pelajarin pass string ini ke handler
	}
	//check if token is expired
	if time.Now().After(RefreshToken.ExpiredAt) {
		return nil, errors.New("refresh token expired")
	}

	//find token by user id
	user, err := s.repo.FindByID(RefreshToken.UserID.String())
	if user == nil {
		return nil, errors.New("user not found") //todo pelajarin pass string ini ke handler
	}

	//generate new access token
	accessToken, err := jwt.GenerateToken(user, s.jwtSecret)
	if err != nil {
		return nil, errors.New("failed to generate token")
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

func (s *userAuthService) Logout(token string) error {
	err := s.refreshTokenRepo.DeleteByToken(token)
	if err != nil {
		return errors.New("failed to delete refresh token")
	}
	return nil
}

// verify email
func (s *userAuthService) VerifyEmail(token string) error {
	verificationToken, err := s.verificationTokenRepo.FindByToken(token)

	if err != nil {
		return errors.New("failed to find verification token")
	}

	if verificationToken == nil {
		return errors.New("verification token not found")
	}

	if time.Now().After(verificationToken.ExpiredAt) {
		return errors.New("verification token expired")
	}

	//note update boolean user verified is using the dto interface of user repository not userauthservice.
	err = s.repo.UpdateVerified(verificationToken.UserID)

	err = s.verificationTokenRepo.DeleteByToken(token)
	if err != nil {
		return errors.New("failed to delete verification token")
	}

	return nil
}

func (s *userAuthService) ResendVerificationEmail(email string) error {
	return nil
}
