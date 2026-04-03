package handler

import (
	"strings"

	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/oauth"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var validate = validator.New()

type authHandler struct {
	service     contract.UserAuthService
	googleOAuth *oauth.GoogleOAuthService
}

func (h *authHandler) register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		thing := response.ErrBadRequest("body format is invalid")
		return response.Error(ctx, thing, err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	res, apiErr := h.service.Register(&req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusCreated, "Registration Success", res)

}

func (h *authHandler) login(ctx *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("body format is invalid"), err)
	}
	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	res, apiErr := h.service.Login(&req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	// Set refresh token ke HTTP-only cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		MaxAge:   int(time.Until(res.RefreshTokenExpiresAt).Seconds()), // 7 days
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	// Hapus refresh token dari body
	res.RefreshToken = ""

	return response.Success(ctx, fiber.StatusOK, "Login Success", res)
}
func (h *authHandler) refresh(ctx *fiber.Ctx) error {
	// Ambil refresh token dari cookie, bukan body
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Error(ctx, response.ErrUnAuthorized("refresh token not found"), nil)
	}

	res, apiErr := h.service.RefreshToken(refreshToken)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	// Perbarui cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		MaxAge:   int(time.Until(res.RefreshTokenExpiresAt).Seconds()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	res.RefreshToken = ""

	return response.Success(ctx, fiber.StatusOK, "Your session has been refreshed", res)
}

func (h *authHandler) logout(ctx *fiber.Ctx) error {
	// Ambil dari cookie
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return response.Error(ctx, response.ErrBadRequest("refresh token not found"), nil)
	}

	apiErr := h.service.Logout(refreshToken)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	// Hapus cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		HTTPOnly: true,
	})

	return response.Success(ctx, fiber.StatusOK, "Logout Success, see you later! bye", nil)
}

func (h *authHandler) verifyEmail(ctx *fiber.Ctx) error {
	token := strings.TrimSpace(ctx.Query("token"))

	frontendLoginURL := "https://work-able-app.vercel.app/verify"

	if token == "" {

		return ctx.Redirect(frontendLoginURL+"?error=token_empty", fiber.StatusTemporaryRedirect)
	}

	//service check token match
	apiErr := h.service.VerifyEmail(token)
	if apiErr != nil {
		return ctx.Redirect(frontendLoginURL+"?error=verification_failed", fiber.StatusTemporaryRedirect)
	}

	return ctx.Redirect(frontendLoginURL+"?verified=true", fiber.StatusTemporaryRedirect)
}

func (h *authHandler) googleLogin(ctx *fiber.Ctx) error {
	state := uuid.New().String()

	//save cookie
	ctx.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 minutes
		HTTPOnly: true,
	})

	// redirect to Google
	url := h.googleOAuth.GetAuthURL(state)
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (h *authHandler) googleCallback(ctx *fiber.Ctx) error {
	state := ctx.Query("state")
	cookieState := ctx.Cookies("oauth_state")
	if !h.googleOAuth.VerifyState(state, cookieState) {
		return response.Error(ctx, response.ErrBadRequest("invalid oauth state"), nil)
	}

	code := ctx.Query("code")
	if code == "" {
		return response.Error(ctx, response.ErrBadRequest("authorization code is required"), nil)
	}

	token, err := h.googleOAuth.ExchangeCode(ctx.Context(), code)
	if err != nil {
		return response.Error(ctx, response.ErrInternal("failed to exchange code"), err)
	}

	userInfo, err := h.googleOAuth.GetUserInfo(ctx.Context(), token)
	if err != nil {
		return response.Error(ctx, response.ErrInternal("failed to get user info"), err)
	}

	res, apiErr := h.service.GoogleAuth(&dto.GoogleAuthRequest{
		Email:   userInfo.Email,
		Name:    userInfo.Name,
		Picture: userInfo.Picture,
	})
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	// Set refresh token ke cookie juga untuk Google login
	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    res.RefreshToken,
		MaxAge:   int(time.Until(res.RefreshTokenExpiresAt).Seconds()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Non",
	})

	res.RefreshToken = ""

	return response.Success(ctx, fiber.StatusOK, "Google login successful", res)
}

func (h *authHandler) resendVerification(ctx *fiber.Ctx) error {
	var req dto.ResendVerificationRequest

	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid request format"), err)
	}

	if apiErr := h.service.ResendVerificationEmail(req.Email); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Verification email has been resent successfully", nil)
}

// Tambahin di struct handler lu

func (h *authHandler) ForgotPassword(ctx *fiber.Ctx) error {
	var req dto.ForgotPasswordRequest

	// Parsing body request dari FE
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Invalid Request Format"), nil)
	}

	apiErr := h.service.ForgotPassword(req.Email)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "if email registered, reset password will be sent to your email", nil)
}

func (h *authHandler) ResetPassword(ctx *fiber.Ctx) error {
	var req dto.ResetPasswordRequest

	// Parsing body from FE
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Invalid Request Format"), nil)
	}

	apiErr := h.service.ResetPassword(req.Token, req.NewPassword)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "password succesfully changed.", nil)
}
