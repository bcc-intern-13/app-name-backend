package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type authHandler struct {
	service dto.UserAuthService
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

	res, err := h.service.Register(&req)
	if err != nil {
		switch err.Error() {
		case "user is already registered":
			return response.Error(ctx, response.ErrConflict(err.Error()), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
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

	res, err := h.service.Login(&req)
	if err != nil {
		switch err.Error() {
		case "user not found", "Wrong password, please try again":
			return response.Error(ctx, response.ErrUnAuthorized("invalid email or password"), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}

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
	res, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err.Error() {
		case "refresh token expired", "refresh token not found", "user not found":
			return response.Error(ctx, response.ErrUnAuthorized("session expired, please login again"), err)
		default:
			return response.Error(ctx, response.ErrInternal(err.Error()), err)
		}
	}

	// return Successs
	return response.Success(ctx, fiber.StatusOK, "Your session has been refreshed", res)
	//todo  dan pahmin ini

}

func (h *authHandler) logout(ctx *fiber.Ctx) error {
	//ambil dto refresh
	var req dto.RefreshRequest

	if err := ctx.BodyParser(&req); err != nil {
		thing := response.ErrBadRequest("invalid body format")
		return response.Error(ctx, thing, err)
	}

	if err := validate.Struct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	// panggil logout
	res, err := h, h.service.Logout(req.RefreshToken)
	if err != nil {
		return response.Error(ctx, response.ErrInternal(err.Error()), err)
	}

	// return
	return response.Success(ctx, fiber.StatusOK, "Logout Success, see you later! bye", res)
}

//todo selesain authnya dulu , basic crud fiturnya
//note wajib pahamin hal hal yang dilakuin
//todo domain driven architecture
//todo penugasan erd
//todo refresh token
