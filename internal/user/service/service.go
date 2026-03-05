package service

import (
	"errors"

	jwt "github.com/bcc-intern-13/app-name-backend/internal/middleware/jwt"
	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userAuthService struct {
	repo      dto.UserRepository
	jwtSecret string
}

func NewUserAuthService(repo dto.UserRepository, jwtSecret string) dto.UserAuthService {
	return &userAuthService{
		repo:      repo,
		jwtSecret: jwtSecret,
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

	return &dto.LoginResponse{
		AccessToken: accessToken,
		User: dto.UserData{
			ID:    user.ID.String(),
			Email: user.Email,
		},
	}, nil

}
