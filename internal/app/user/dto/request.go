package dto

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=2"`
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email"     validate:"required,email"`
	Password string `json:"password"  validate:"required,min=8"`
}
type LoginRequest struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password"  validate:"required,min=8"`
	RememberMe bool   `json:"remember_me"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type GoogleAuthRequest struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}
