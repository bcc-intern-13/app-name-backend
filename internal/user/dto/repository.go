package dto

import (
	"github.com/bcc-intern-13/app-name-backend/internal/user/entity"
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
}

type VerificationTokenRepository interface {
	Create(token *entity.VerificationToken) error
	FindByToken(token string) (*entity.VerificationToken, error)
	DeleteByToken(token string) error
}
