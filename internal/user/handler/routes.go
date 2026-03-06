package handler

import (
	"github.com/bcc-intern-13/app-name-backend/internal/middleware"
	"github.com/bcc-intern-13/app-name-backend/internal/user/dto"
	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, userService dto.UserAuthService, jwtSecret string) {
	// Public
	auth := app.Group("/auth")
	authH := &authHandler{service: userService}
	auth.Post("/register", authH.register)
	auth.Post("/login", authH.login)
	auth.Post("/refresh", authH.refresh) // refresh)
	auth.Post("/logout", authH.logout)

	// Protected
	userH := &userHandler{service: userService}
	users := app.Group("/users", middleware.JWTProtected(jwtSecret))
	users.Get("/me", userH.getMe)
}
