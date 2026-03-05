package middleware

import (
	"strings"

	"github.com/bcc-intern-13/app-name-backend/internal/middleware/jwt"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func JWTProtected(secret string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return response.Error(ctx, response.ErrUnAuthorized("missing token"), nil)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := jwt.ParseToken(tokenStr, secret)
		if err != nil {
			return response.Error(ctx, response.ErrUnAuthorized("invalid or expired token"), err)
		}

		//todo pahamin kodingan ini
		ctx.Locals("userID", claims.UserID)
		ctx.Locals("email", claims.Email)

		return ctx.Next()
	}
}
