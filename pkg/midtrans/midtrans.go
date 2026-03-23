package midtrans

import (
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransService struct {
	client      snap.Client
	serverKey   string
	environment midtrans.EnvironmentType
}

func NewMidtransService(serverKey string, isProduction bool) *MidtransService {
	env := midtrans.Sandbox
	if isProduction {
		env = midtrans.Production
	}

	var client snap.Client
	client.New(serverKey, env)

	return &MidtransService{
		client:      client,
		serverKey:   serverKey,
		environment: env,
	}
}

// CreateTransaction → kirim request ke Midtrans, dapat payment_url
func (m *MidtransService) CreateTransaction(orderID string, amount int64, userEmail, userName string) (string, error) {
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			Email: userEmail,
			FName: userName,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "WORKABLE_PREMIUM",
				Name:  "WorkAble Premium - CV AI Access",
				Price: amount,
				Qty:   1,
			},
		},
	}

	snapResp, err := m.client.CreateTransaction(req)
	if err != nil {
		return "", fmt.Errorf("midtrans create transaction: %w", err)
	}

	return snapResp.RedirectURL, nil
}

// VerifyWebhook → verify signature dari Midtrans notification
// Midtrans signature = SHA512(order_id + status_code + gross_amount + server_key)
func (m *MidtransService) VerifyWebhook(orderID, statusCode, grossAmount, signatureKey string) bool {
	expectedSignature := midtrans.GetSignatureKey(orderID, statusCode, grossAmount, m.serverKey)
	return expectedSignature == signatureKey
}
