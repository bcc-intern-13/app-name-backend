package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/domain/dto"
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "invalid request body",
			"detail":  "The data you sent is incorrect, please check again",
			"status":  400,
		})
	}

	if err := validate.Struct(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	res, err := h.service.Register(&req)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Registration successful, please check your email for verification",
		"data":    res,
	})
}

//todo error custom, response error, 2 macem buat end user dan dev team
//todo selesain authnya dulu , basic crud fiturnya
//note wajib pahamin hal hal yang dilakuin
//todo domain driven architecture
//todo penugasan erd
//todo refresh token
