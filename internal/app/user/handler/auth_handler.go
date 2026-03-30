package handler

import (
	"strings"

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
		thing := response.ErrBadRequest("body format is invalid")
		return response.Error(ctx, thing, err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	res, apiErr := h.service.Login(&req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}
	return response.Success(ctx, fiber.StatusCreated, "Login Success", res)

}

func (h *authHandler) refresh(ctx *fiber.Ctx) error {
	// get refresh token base on refresh request
	var req dto.RefreshRequest

	if err := ctx.BodyParser(&req); err != nil {
		thing := response.ErrBadRequest("body format is invalid")
		return response.Error(ctx, thing, err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	//refresh token
	res, apiErr := h.service.RefreshToken(req.RefreshToken)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	// return Successs
	return response.Success(ctx, fiber.StatusOK, "Your session has been refreshed", res)

}

func (h *authHandler) logout(ctx *fiber.Ctx) error {
	var req dto.RefreshRequest

	if err := ctx.BodyParser(&req); err != nil {
		thing := response.ErrBadRequest("invalid body format")
		return response.Error(ctx, thing, err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	apiErr := h.service.Logout(req.RefreshToken)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	// return
	return response.Success(ctx, fiber.StatusOK, "Logout Success, see you later! bye", nil)
}

func (h *authHandler) verifyEmail(ctx *fiber.Ctx) error {

	token := strings.TrimSpace(ctx.Query("token"))

	if token == "" {
		return response.Error(ctx, response.ErrBadRequest("token is required"), nil)
	}

	apiErr := h.service.VerifyEmail(token)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}
	return response.Success(ctx, fiber.StatusOK, "Email verification successful, you can now login", nil)

}

func (h *authHandler) googleLogin(ctx *fiber.Ctx) error {
	// generate random state untuk prevent CSRF
	state := uuid.New().String()

	// simpan state ke cookie sementara
	ctx.Cookie(&fiber.Cookie{
		Name:     "oauth_state",
		Value:    state,
		MaxAge:   300, // 5 menit
		HTTPOnly: true,
	})

	// redirect ke Google
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

	return response.Success(ctx, fiber.StatusOK, "Google login successful", res)
}
