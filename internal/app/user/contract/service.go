package contract

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/dto"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
)

type UserAuthService interface {
	Register(req *dto.RegisterRequest) (*entity.User, *response.APIError)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, *response.APIError)
	//using refresh token, to refresh
	RefreshToken(token string) (*dto.LoginResponse, *response.APIError)
	Logout(token string) *response.APIError

	//verification gmial
	VerifyEmail(token string) *response.APIError
	ResendVerificationEmail(email string) *response.APIError
}
