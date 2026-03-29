package middleware

import (
	"time"

	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func PremiumRequired(userRepo contract.UserRepository) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userIDStr := ctx.Locals("userID").(string)
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return response.Error(ctx, response.ErrBadRequest("invalid user id"), err)
		}

		user, err := userRepo.FindByID(userID.String())
		if err != nil || user == nil {
			return response.Error(ctx, response.ErrInternal("failed to get user"), err)
		}

		// check is_premium + expired
		if !user.IsPremium {
			return response.Error(ctx, response.ErrForbidden("this feature is for premium users only"), nil)
		}

		// check if user's premium has expired
		if user.PremiumExpiresAt != nil && user.PremiumExpiresAt.Before(time.Now()) {
			// auto reset is_premium ke false
			user.IsPremium = false
			user.PremiumExpiresAt = nil
			// update DB fire and forget
			go userRepo.UpdatePremiumStatus(userID, false, nil)
			return response.Error(ctx, response.ErrForbidden("your premium has expired, please renew"), nil)
		}

		return ctx.Next()
	}
}
