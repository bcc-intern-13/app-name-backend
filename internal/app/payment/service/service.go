package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/bcc-intern-13/app-name-backend/internal/app/payment/contract"
	"github.com/bcc-intern-13/app-name-backend/internal/app/payment/dto"
	"github.com/bcc-intern-13/app-name-backend/internal/app/payment/entity"
	userContract "github.com/bcc-intern-13/app-name-backend/internal/app/user/contract"
	"github.com/bcc-intern-13/app-name-backend/pkg/response"
	"github.com/bcc-intern-13/app-name-backend/pkg/xendit"
	"github.com/google/uuid"
)

const premiumPrice float64 = 99000 // Rp 99.000

type paymentService struct {
	orderRepo    contract.OrderRepository
	userRepo     userContract.UserRepository
	xendit       *xendit.XenditService
	webhookToken string
}

func NewPaymentService(
	orderRepo contract.OrderRepository,
	userRepo userContract.UserRepository,
	xenditSvc *xendit.XenditService,
	webhookToken string,
) contract.PaymentService {
	return &paymentService{
		orderRepo:    orderRepo,
		userRepo:     userRepo,
		xendit:       xenditSvc,
		webhookToken: webhookToken,
	}
}

func (s *paymentService) CreateOrder(ctx context.Context, userID uuid.UUID) (*dto.CreateOrderResponse, *response.APIError) {
	user, err := s.userRepo.FindByID(userID.String())
	if err != nil {
		slog.Error("failed to get user", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get user")
	}
	if user == nil {
		return nil, response.ErrNotFound("user not found")
	}
	if user.IsPremium {
		return nil, response.ErrConflict("user is already premium")
	}

	pending, err := s.orderRepo.FindPendingByUserID(userID)
	if err != nil {
		return nil, response.ErrInternal("failed to check pending order")
	}
	if pending != nil {
		return &dto.CreateOrderResponse{
			OrderID:    pending.ID,
			PaymentURL: pending.PaymentURL,
			Amount:     pending.Amount,
		}, nil
	}

	// make external ID unique for xendit, format: WORKABLE-{userID}-{timestamp}å
	externalID := fmt.Sprintf("WORKABLE-%s-%d", userID.String()[:8], time.Now().Unix())

	// kirim ke Xendit, dapat invoice_url
	invoiceURL, err := s.xendit.CreateInvoice(
		ctx,
		externalID,
		premiumPrice,
		user.Email,
		user.Nama,
	)
	if err != nil {
		slog.Error("failed to create xendit invoice", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to create payment")
	}

	// save order to Database
	order := &entity.Order{
		ID:               uuid.New(),
		UserID:           userID,
		XenditExternalID: externalID,
		Amount:           int64(premiumPrice),
		Status:           "PENDING",
		PaymentURL:       invoiceURL,
	}
	if err := s.orderRepo.Create(order); err != nil {
		slog.Error("failed to create order", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to save order")
	}

	return &dto.CreateOrderResponse{
		OrderID:    order.ID,
		PaymentURL: invoiceURL,
		Amount:     int64(premiumPrice),
	}, nil
}

func (s *paymentService) HandleWebhook(ctx context.Context, req *dto.WebhookRequest) *response.APIError {
	// verify callback token
	if !s.xendit.VerifyWebhook(req.CallbackToken, s.webhookToken) {
		slog.Error("invalid xendit webhook token", "externalID", req.ExternalID)
		return response.ErrUnAuthorized("invalid callback token")
	}

	// find order by xendit's external_id
	order, err := s.orderRepo.FindByXenditExternalID(req.ExternalID)
	if err != nil {
		slog.Error("failed to find order", "error", err, "externalID", req.ExternalID)
		return response.ErrInternal("failed to find order")
	}
	if order == nil {
		return response.ErrNotFound("order not found")
	}

	switch req.Status {
	case "PAID":
		order.Status = "PAID"
		order.PaymentType = req.PaymentMethod

		// upgrade user to premium
		user, err := s.userRepo.FindByID(order.UserID.String())
		if err != nil || user == nil {
			slog.Error("failed to find user for premium upgrade", "error", err, "userID", order.UserID)
			return response.ErrInternal("failed to upgrade user")
		}
		user.IsPremium = true
		if err := s.userRepo.UpdateIsPremium(user.ID); err != nil {
			slog.Error("failed to upgrade user to premium", "error", err, "userID", order.UserID)
			return response.ErrInternal("failed to upgrade user to premium")
		}
		slog.Info("user upgraded to premium", "userID", order.UserID)

	case "EXPIRED":
		order.Status = "EXPIRED"
	}

	if err := s.orderRepo.Update(order); err != nil {
		slog.Error("failed to update order", "error", err, "externalID", req.ExternalID)
		return response.ErrInternal("failed to update order")
	}

	return nil
}

func (s *paymentService) GetOrderHistory(ctx context.Context, userID uuid.UUID) ([]dto.OrderResponse, *response.APIError) {
	orders, err := s.orderRepo.FindByUserID(userID)
	if err != nil {
		slog.Error("failed to get order history", "error", err, "userID", userID)
		return nil, response.ErrInternal("failed to get order history")
	}

	var result []dto.OrderResponse
	for _, o := range orders {
		result = append(result, dto.OrderResponse{
			ID:               o.ID,
			XenditExternalID: o.XenditExternalID,
			Amount:           o.Amount,
			Status:           o.Status,
			PaymentType:      o.PaymentType,
			CreatedAt:        o.CreatedAt,
		})
	}

	return result, nil
}
