package handler

//note temporal hanlder for user_handler.
import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userHandler struct {
	service contract.UserAuthService
}

func (h *userHandler) getMe(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(string)
	email := ctx.Locals("email").(string)

	return response.Success(ctx, fiber.StatusOK, "success", fiber.Map{
		"user_id": userID,
		"email":   email,
	})
}

func (h *userHandler) UploadAvatar(ctx *fiber.Ctx) error {
	// 1. Ambil ID dari JWT Token (Pastikan key-nya "userID" sesuai dengan middleware JWT lu)
	userIDStr := ctx.Locals("userID").(string)
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrUnAuthorized("invalid user token"), err)
	}

	// 2. Tangkap file gambarnya (Zildjian harus kirim pakai key "avatar")
	file, err := ctx.FormFile("avatar")
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("avatar file is required"), err)
	}

	// 3. Eksekusi ke Service
	res, apiErr := h.service.UploadAvatar(ctx.Context(), userID, file)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Avatar uploaded successfully", res)
}
