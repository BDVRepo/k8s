package paymentsgrpc

import (
	"context"
	"log"

	paymentspb "github.com/bdv/gprs/internal/gen/paymentspb"
	paymentsservice "github.com/bdv/gprs/internal/methods/handlers/payments/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	paymentspb.UnimplementedPaymentsServiceServer
	service *paymentsservice.Service
}

func NewGRPCServer(service *paymentsservice.Service) *GRPCServer {
	return &GRPCServer{service: service}
}

func (s *GRPCServer) AuthorizePayment(ctx context.Context, req *paymentspb.AuthorizePaymentRequest) (*paymentspb.AuthorizePaymentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	log.Printf("[payments.grpc] AuthorizePayment request order_id=%s amount=%.2f %s",
		req.GetOrderId(), req.GetAmount(), req.GetCurrency())

	payment, err := s.service.AuthorizePayment(ctx, paymentsservice.AuthorizeCommand{
		OrderID:  req.GetOrderId(),
		Amount:   req.GetAmount(),
		Currency: req.GetCurrency(),
	})
	if err != nil {
		log.Printf("[payments.grpc] AuthorizePayment failed order_id=%s err=%v", req.GetOrderId(), err)
		return nil, status.Errorf(codes.InvalidArgument, "authorize payment: %v", err)
	}

	log.Printf("[payments.grpc] AuthorizePayment success payment_id=%s status=%s", payment.ID, payment.Status)

	return &paymentspb.AuthorizePaymentResponse{
		PaymentId: payment.ID,
		Status:    string(payment.Status),
		Message:   payment.Message,
	}, nil
}

func (s *GRPCServer) GetPaymentStatus(ctx context.Context, req *paymentspb.GetPaymentStatusRequest) (*paymentspb.GetPaymentStatusResponse, error) {
	if req == nil || req.GetPaymentId() == "" {
		return nil, status.Error(codes.InvalidArgument, "payment_id is required")
	}

	log.Printf("[payments.grpc] GetPaymentStatus request payment_id=%s", req.GetPaymentId())

	payment, err := s.service.GetPayment(ctx, req.GetPaymentId())
	if err != nil {
		log.Printf("[payments.grpc] GetPaymentStatus failed payment_id=%s err=%v", req.GetPaymentId(), err)
		return nil, status.Errorf(codes.NotFound, "get payment: %v", err)
	}

	log.Printf("[payments.grpc] GetPaymentStatus success payment_id=%s status=%s", payment.ID, payment.Status)

	return &paymentspb.GetPaymentStatusResponse{
		PaymentId: payment.ID,
		Status:    string(payment.Status),
		Message:   payment.Message,
	}, nil
}
