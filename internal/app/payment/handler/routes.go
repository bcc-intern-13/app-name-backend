package handler

import (
	"github.com/bcc-intern-13/WorkAble-backend/internal/app/payment/contract"
	"github.com/bcc-intern-13/WorkAble-backend/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func RegisterPaymentRoutes(router fiber.Router, service contract.PaymentService, jwtSecret string) {
	h := &paymentHandler{service: service}

	payment := router.Group("api/v1/payment")
	//xendit will hit this endpoint to update oder table
	payment.Post("/webhook", h.handleWebhook)

	payment.Post("/create-order", middleware.JWTProtected(jwtSecret), h.createOrder)
	payment.Get("/history", middleware.JWTProtected(jwtSecret), h.getOrderHistory)
}
