package handler

import (
	"strings"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/dto"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type authHandler struct {
	service contract.UserAuthService
}

func (h *authHandler) register(ctx *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := ctx.BodyParser(&req); err != nil {
		thing := response.ErrBadRequest("body format is invalid") //note mencoba belajar parse ke dalam parameter langsung dan di inferensikan ke variabel dulu
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
		thing := response.ErrBadRequest("body format is invalid") //note mencoba belajar parse ke dalam parameter langsung dan di inferensikan ke variabel dulu
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
		thing := response.ErrBadRequest("body format is invalid") //note mencoba belajar parse ke dalam parameter langsung dan di inferensikan ke variabel dulu
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
	//todo  dan pahmin ini

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

//note wajib pahamin hal hal yang dilakuin
//todo penugasan erd
