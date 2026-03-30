package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/bcc-intern-13/WorkAble-backend/pkg/oauth"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userService contract.UserAuthService, jwtSecret string, googleOAuth *oauth.GoogleOAuthService) {
	// Public
	auth := app.Group("api/v1/auth")
	authH := &authHandler{
		service:     userService,
		googleOAuth: googleOAuth,
	}
	auth.Post("/register", authH.register)
	auth.Post("/login", authH.login)
	auth.Post("/logout", authH.logout)

	//refresh acces token usign refresh token
	auth.Post("/refresh-token", authH.refresh)

	//verify gmail
	auth.Get("/verify", authH.verifyEmail)

	//google login
	auth.Get("/google", authH.googleLogin)
	auth.Get("/google/callback", authH.googleCallback)

	// Protected
	userH := &userHandler{service: userService}
	users := app.Group("/users", middleware.JWTProtected(jwtSecret))
	users.Get("/me", userH.getMe)
}
