package ordersrepository

import (
	"context"
	"errors"
	"log"
	"sync"

	ordermodel "github.com/bdv/gprs/internal/models/orders"
)

type OrderRepository interface {
	Save(ctx context.Context, order *ordermodel.Order) error
	GetByID(ctx context.Context, id string) (*ordermodel.Order, error)
}

type InMemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]*ordermodel.Order
}

func NewInMemoryRepository() OrderRepository {
	return &InMemoryRepository{orders: make(map[string]*ordermodel.Order)}
}

func (r *InMemoryRepository) Save(_ context.Context, order *ordermodel.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[order.ID] = order
	log.Printf("[orders.repo] stored order_id=%s status=%s payment_id=%s payment_status=%s",
		order.ID, order.Status, order.PaymentID, order.PaymentStatus)
	return nil
}

func (r *InMemoryRepository) GetByID(_ context.Context, id string) (*ordermodel.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	order, ok := r.orders[id]
	if !ok {
		log.Printf("[orders.repo] order not found order_id=%s", id)
		return nil, errors.New("order not found")
	}
	log.Printf("[orders.repo] fetched order_id=%s status=%s payment_status=%s",
		order.ID, order.Status, order.PaymentStatus)
	return order, nil
}
