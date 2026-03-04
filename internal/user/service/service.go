package service

import (
	"errors"

	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userAuthService struct {
	repo dto.UserRepository
}

func NewUserAuthService(repo dto.UserRepository) dto.UserAuthService {
	return &userAuthService{repo: repo}
}

func (s *userAuthService) Register(req *dto.RegisterRequest) (*entity.User, error) {
	existing, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email is already registered")
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
