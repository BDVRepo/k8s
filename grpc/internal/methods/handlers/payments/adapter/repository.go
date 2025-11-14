package paymentsrepository

import (
	"context"
	"errors"
	"log"
	"sync"

	paymentmodel "github.com/bdv/gprs/internal/models/payments"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment *paymentmodel.Payment) error
	GetByID(ctx context.Context, id string) (*paymentmodel.Payment, error)
}

type InMemoryRepository struct {
	mu       sync.RWMutex
	payments map[string]*paymentmodel.Payment
}

func NewInMemoryRepository() PaymentRepository {
	return &InMemoryRepository{payments: make(map[string]*paymentmodel.Payment)}
}

func (r *InMemoryRepository) Save(_ context.Context, payment *paymentmodel.Payment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.payments[payment.ID] = payment
	log.Printf("[payments.repo] stored payment_id=%s order_id=%s status=%s",
		payment.ID, payment.OrderID, payment.Status)
	return nil
}

func (r *InMemoryRepository) GetByID(_ context.Context, id string) (*paymentmodel.Payment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	payment, ok := r.payments[id]
	if !ok {
		log.Printf("[payments.repo] payment not found payment_id=%s", id)
		return nil, errors.New("payment not found")
	}
	log.Printf("[payments.repo] fetched payment_id=%s status=%s", payment.ID, payment.Status)
	return payment, nil
}
