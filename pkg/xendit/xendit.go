package xendit

import (
	"context"
	"fmt"
	"log"

	"github.com/xendit/xendit-go/v6"
	"github.com/xendit/xendit-go/v6/invoice"
)

type XenditService struct {
	client    *xendit.APIClient
	secretKey string
}

func NewXenditService(secretKey string) *XenditService {
	client := xendit.NewClient(secretKey)
	return &XenditService{
		client:    client,
		secretKey: secretKey,
	}
}

// CreateInvoice → buat invoice di Xendit, return invoice_url
func (x *XenditService) CreateInvoice(ctx context.Context, orderID string, amount float64, userEmail, userName string) (string, error) {
	log.Printf("Xendit request - orderID: %s, amount: %f, email: %s, name: %s", orderID, amount, userEmail, userName)

	description := "WorkAble Premium - CV AI Access"

	req := invoice.NewCreateInvoiceRequest(orderID, amount)
	req.PayerEmail = &userEmail
	req.Description = &description

	customer := invoice.NewCustomerObjectWithDefaults()
	customer.Email = *invoice.NewNullableString(&userEmail)
	customer.GivenNames = *invoice.NewNullableString(&userName)
	if userName != "" {
		customer.GivenNames = *invoice.NewNullableString(&userName)
	}
	req.Customer = customer

	//test
	fmt.Printf("Xendit request - orderID: %s, amount: %f, email: %s, name: %s\n", orderID, amount, userEmail, userName)

	resp, _, sdkErr := x.client.InvoiceApi.CreateInvoiceExecute(
		x.client.InvoiceApi.CreateInvoice(ctx).CreateInvoiceRequest(*req),
	)
	if sdkErr != nil {
		return "", fmt.Errorf("xendit create invoice: %s", sdkErr.Error())
	}

	return resp.InvoiceUrl, nil
}

// VerifyWebhook → verify callback token dari header x-callback-token
func (x *XenditService) VerifyWebhook(callbackToken, expectedToken string) bool {
	return callbackToken == expectedToken
}
