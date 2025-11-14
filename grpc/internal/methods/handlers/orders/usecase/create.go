package ordersservice

import (
	"context"
	"log"

	grpcclient "github.com/bdv/gprs/internal/infrastructure/grpc"
	ordersrepository "github.com/bdv/gprs/internal/methods/handlers/orders/adapter"
	ordermodel "github.com/bdv/gprs/internal/models/orders"
	"github.com/bdv/gprs/internal/models/shared"
)

type PaymentAuthorizer = grpcclient.PaymentAuthorizer

type IDGenerator interface {
	NewID() string
}

type CreateOrderCommand struct {
	UserID    string
	ProductID string
	Price     shared.Money
}

type Service struct {
	repo     ordersrepository.OrderRepository
	payments PaymentAuthorizer
	idGen    IDGenerator
}

func NewService(repo ordersrepository.OrderRepository, payments PaymentAuthorizer, idGen IDGenerator) *Service {
	return &Service{
		repo:     repo,
		payments: payments,
		idGen:    idGen,
	}
}

func (s *Service) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (*ordermodel.Order, error) {
	orderID := s.idGen.NewID()
	order := ordermodel.New(orderID, cmd.UserID, cmd.ProductID, cmd.Price)

	log.Printf("[orders.service] create start order_id=%s user_id=%s product_id=%s price=%.2f %s",
		orderID, cmd.UserID, cmd.ProductID, cmd.Price.Amount, cmd.Price.Currency)

	paymentID, paymentStatus, err := s.payments.Authorize(ctx, orderID, cmd.Price)
	if err != nil {
		log.Printf("[orders.service] payment authorize failed order_id=%s err=%v", orderID, err)
		return nil, err
	}

	log.Printf("[orders.service] payment authorize success order_id=%s payment_id=%s payment_status=%s",
		orderID, paymentID, paymentStatus)

	order.MarkPayment(paymentID, paymentStatus)

	if err := s.repo.Save(ctx, order); err != nil {
		log.Printf("[orders.service] save failed order_id=%s err=%v", orderID, err)
		return nil, err
	}

	log.Printf("[orders.service] create done order_id=%s status=%s", orderID, order.Status)

	return order, nil
}
