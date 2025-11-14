package ordersgrpc

import (
	"context"
	"log"

	orderspb "github.com/bdv/gprs/internal/gen/orderspb"
	ordersservice "github.com/bdv/gprs/internal/methods/handlers/orders/usecase"
	"github.com/bdv/gprs/internal/models/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	orderspb.UnimplementedOrdersServiceServer
	service *ordersservice.Service
}

func NewGRPCServer(service *ordersservice.Service) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) CreateOrder(ctx context.Context, req *orderspb.CreateOrderRequest) (*orderspb.CreateOrderResponse, error) {
	if req == nil || req.Price == nil {
		return nil, status.Error(codes.InvalidArgument, "price is required")
	}

	log.Printf("[orders.grpc] CreateOrder request received user_id=%s product_id=%s currency=%s amount=%.2f",
		req.GetUserId(), req.GetProductId(), req.GetPrice().GetCurrency(), req.GetPrice().GetAmount())

	cmd := ordersservice.CreateOrderCommand{
		UserID:    req.GetUserId(),
		ProductID: req.GetProductId(),
		Price: shared.Money{
			Currency: req.Price.GetCurrency(),
			Amount:   req.Price.GetAmount(),
		},
	}

	order, err := s.service.CreateOrder(ctx, cmd)
	if err != nil {
		log.Printf("[orders.grpc] CreateOrder failed user_id=%s product_id=%s err=%v", req.GetUserId(), req.GetProductId(), err)
		return nil, status.Errorf(codes.Internal, "create order: %v", err)
	}

	log.Printf("[orders.grpc] CreateOrder success order_id=%s payment_id=%s status=%s",
		order.ID, order.PaymentID, order.Status)

	return &orderspb.CreateOrderResponse{
		OrderId:   order.ID,
		PaymentId: order.PaymentID,
		Status:    string(order.Status),
	}, nil
}

func (s *GRPCServer) GetOrder(ctx context.Context, req *orderspb.GetOrderRequest) (*orderspb.GetOrderResponse, error) {
	if req == nil || req.GetOrderId() == "" {
		return nil, status.Error(codes.InvalidArgument, "order_id is required")
	}

	log.Printf("[orders.grpc] GetOrder request received order_id=%s", req.GetOrderId())

	order, err := s.service.GetOrder(ctx, req.GetOrderId())
	if err != nil {
		log.Printf("[orders.grpc] GetOrder missed order_id=%s err=%v", req.GetOrderId(), err)
		return nil, status.Errorf(codes.NotFound, "get order: %v", err)
	}

	log.Printf("[orders.grpc] GetOrder success order_id=%s status=%s payment_status=%s",
		order.ID, order.Status, order.PaymentStatus)

	return &orderspb.GetOrderResponse{
		Order: &orderspb.OrderSummary{
			OrderId:       order.ID,
			UserId:        order.UserID,
			ProductId:     order.ProductID,
			Status:        string(order.Status),
			PaymentStatus: order.PaymentStatus,
			PaymentId:     order.PaymentID,
			Price: &orderspb.Money{
				Currency: order.Price.Currency,
				Amount:   order.Price.Amount,
			},
		},
	}, nil
}
