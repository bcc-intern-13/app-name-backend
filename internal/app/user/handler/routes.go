package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userService contract.UserAuthService, jwtSecret string) {
	// Public
	auth := app.Group("/auth")
	authH := &authHandler{service: userService}
	auth.Post("/register", authH.register)
	auth.Post("/login", authH.login)
	auth.Post("/refresh-token", authH.refresh) // refresh)
	auth.Post("/logout", authH.logout)

	auth.Get("/verify", authH.verifyEmail)

	// Protected
	userH := &userHandler{service: userService}
	users := app.Group("/users", middleware.JWTProtected(jwtSecret))
	users.Get("/me", userH.getMe)
}
