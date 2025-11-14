package paymentsservice

import (
	"context"
	"log"

	paymentmodel "github.com/bdv/gprs/internal/models/payments"
)

func (s *Service) GetPayment(ctx context.Context, paymentID string) (*paymentmodel.Payment, error) {
	log.Printf("[payments.service] fetch payment_id=%s", paymentID)
	return s.repo.GetByID(ctx, paymentID)
}
