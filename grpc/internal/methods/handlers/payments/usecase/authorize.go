package paymentsservice

import (
	"context"
	"errors"
	"log"

	paymentsrepository "github.com/bdv/gprs/internal/methods/handlers/payments/adapter"
	paymentmodel "github.com/bdv/gprs/internal/models/payments"
)

type IDGenerator interface {
	NewID() string
}

type AuthorizeCommand struct {
	OrderID  string
	Amount   float64
	Currency string
}

type Service struct {
	repo  paymentsrepository.PaymentRepository
	idGen IDGenerator
}

func NewService(repo paymentsrepository.PaymentRepository, idGen IDGenerator) *Service {
	return &Service{
		repo:  repo,
		idGen: idGen,
	}
}

func (s *Service) AuthorizePayment(ctx context.Context, cmd AuthorizeCommand) (*paymentmodel.Payment, error) {
	if cmd.Amount <= 0 {
		log.Printf("[payments.service] invalid amount order_id=%s amount=%.2f", cmd.OrderID, cmd.Amount)
		return nil, errors.New("amount must be positive")
	}

	paymentID := s.idGen.NewID()
	payment := paymentmodel.New(paymentID, cmd.OrderID, cmd.Amount, cmd.Currency, paymentmodel.StatusAuthorized, "payment authorized")

	log.Printf("[payments.service] authorize start payment_id=%s order_id=%s amount=%.2f %s",
		paymentID, cmd.OrderID, cmd.Amount, cmd.Currency)

	if err := s.repo.Save(ctx, payment); err != nil {
		log.Printf("[payments.service] repo save failed payment_id=%s err=%v", paymentID, err)
		return nil, err
	}

	log.Printf("[payments.service] authorize success payment_id=%s status=%s", paymentID, payment.Status)

	return payment, nil
}
