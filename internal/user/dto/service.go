package dto

import "github.com/bcc-intern-13/app-name-backend/internal/user/entity"

type UserAuthService interface {
	Register(req *RegisterRequest) (*entity.User, error)
	Login(req *LoginRequest) (*LoginResponse, error)
	//using refresh token, to refresh
	RefreshToken(token string) (*LoginResponse, error)
	Logout(token string) error

	//verification gmial
	VerifyEmail(token string) error
	ResendVerificationEmail(email string) error
}
