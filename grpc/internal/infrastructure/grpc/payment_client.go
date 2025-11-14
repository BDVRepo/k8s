package grpc

import (
	"context"
	"log"

	"github.com/bdv/gprs/internal/models/shared"
	paymentspb "github.com/bdv/gprs/internal/gen/paymentspb"
)

type PaymentAuthorizer interface {
	Authorize(ctx context.Context, orderID string, amount shared.Money) (paymentID string, status string, err error)
}

type GRPCPaymentClient struct {
	client paymentspb.PaymentsServiceClient
}

func NewGRPCPaymentClient(client paymentspb.PaymentsServiceClient) PaymentAuthorizer {
	return &GRPCPaymentClient{client: client}
}

func (c *GRPCPaymentClient) Authorize(ctx context.Context, orderID string, amount shared.Money) (string, string, error) {
	log.Printf("[orders.payment-client] authorize request order_id=%s amount=%.2f %s",
		orderID, amount.Amount, amount.Currency)

	resp, err := c.client.AuthorizePayment(ctx, &paymentspb.AuthorizePaymentRequest{
		OrderId:  orderID,
		Amount:   amount.Amount,
		Currency: amount.Currency,
	})
	if err != nil {
		log.Printf("[orders.payment-client] authorize failed order_id=%s err=%v", orderID, err)
		return "", "", err
	}

	log.Printf("[orders.payment-client] authorize success order_id=%s payment_id=%s status=%s",
		orderID, resp.GetPaymentId(), resp.GetStatus())

	return resp.GetPaymentId(), resp.GetStatus(), nil
}

