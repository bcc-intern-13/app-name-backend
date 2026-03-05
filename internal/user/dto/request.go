package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}
type UserAuthService interface {
	Register(req *RegisterRequest) (*entity.User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
}

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
}
