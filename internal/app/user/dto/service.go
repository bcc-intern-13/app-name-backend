package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
)

type UserAuthService interface {
	Register(req *RegisterRequest) (*entity.User, *response.APIError)
	Login(req *LoginRequest) (*LoginResponse, *response.APIError)
	//using refresh token, to refresh
	RefreshToken(token string) (*LoginResponse, *response.APIError)
	Logout(token string) *response.APIError

	//verification gmial
	VerifyEmail(token string) *response.APIError
	ResendVerificationEmail(email string) *response.APIError
}
