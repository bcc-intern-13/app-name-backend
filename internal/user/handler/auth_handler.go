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

func NewAuthHandler(app *fiber.App, u dto.UserAuthService) {
	handler := &authHandler{service: u}

	auth := app.Group("/auth")
	auth.Post("/register", handler.register)
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

//todo error custom, response error, 2 macem buat end user dan dev team
//todo selesain authnya dulu , basic crud fiturnya
//note wajib pahamin hal hal yang dilakuin
//todo domain driven architecture
//todo penugasan erd
//todo refresh token
