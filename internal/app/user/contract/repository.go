package contract

import (
	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/entity"
	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(token *entity.RefreshToken) error
	FindByToken(token string) (*entity.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uuid.UUID) error
}

type UserRepository interface {
	FindByEmail(email string) (*entity.User, error)
	Create(user *entity.User) error
	FindByID(id string) (*entity.User, error)
	UpdateVerified(userID uuid.UUID) error
	UpdateOnboardingCompleted(userID uuid.UUID) error
	UpdateIsPremium(userID uuid.UUID) error
	UpdatePremiumStatus(userID uuid.UUID, isPremium bool, expiresAt *time.Time) error
	Update(user *entity.User) error
}

type VerificationTokenRepository interface {
	Create(token *entity.VerificationToken) error
	FindByToken(token string) (*entity.VerificationToken, error)
	DeleteByToken(token string) error
}
