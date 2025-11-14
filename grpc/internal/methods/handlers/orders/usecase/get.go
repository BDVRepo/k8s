package ordersservice

import (
	"context"
	"log"

	ordermodel "github.com/bdv/gprs/internal/models/orders"
)

func (s *Service) GetOrder(ctx context.Context, orderID string) (*ordermodel.Order, error) {
	log.Printf("[orders.service] fetch order_id=%s", orderID)
	return s.repo.GetByID(ctx, orderID)
}

