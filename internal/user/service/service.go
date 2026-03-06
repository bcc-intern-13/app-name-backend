package service

import (
	"errors"
	"time"

	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	jwt "github.com/bcc-intern-13/app-name-backend/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userAuthService struct {
	repo             dto.UserRepository
	refreshTokenRepo dto.RefreshTokenRepository
	jwtSecret        string
}

func NewUserAuthService(repo dto.UserRepository, jwtSecret string, refreshTokenRepo dto.RefreshTokenRepository) dto.UserAuthService {
	return &userAuthService{
		repo:             repo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        jwtSecret,
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
