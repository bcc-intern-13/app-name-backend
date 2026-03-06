package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}

type RefreshTokenRepository interface {
	Create(token *entity.RefreshToken) error
	FindByToken(token string) (*entity.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uuid.UUID) error
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}
type UserAuthService interface {
	Register(req *RegisterRequest) (*entity.User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	//using refresh token, to refresh
	RefreshToken(token string) (*LoginResponse, error)
	Logout(token string) error

	//todo : kerjain ini buat refresh token., implementasiin juga.
	//todo : verify gmail

}

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	FindByID(id string) (*entity.User, error)
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
