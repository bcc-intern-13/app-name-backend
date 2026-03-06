package handler

//note temporal hanlder for user_handler.
import (
	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type userHandler struct {
	service dto.UserAuthService
}

func (h *userHandler) getMe(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	email := ctx.Locals("email").(string)

	return response.Success(ctx, fiber.StatusOK, "success", fiber.Map{
		"user_id": userID,
		"email":   email,
	})
}
