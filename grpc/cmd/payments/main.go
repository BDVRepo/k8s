package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentspb "github.com/bdv/gprs/internal/gen/paymentspb"
	paymentsrepository "github.com/bdv/gprs/internal/methods/handlers/payments/adapter"
	paymentsgrpc "github.com/bdv/gprs/internal/methods/handlers/payments/grpc"
	paymentsservice "github.com/bdv/gprs/internal/methods/handlers/payments/usecase"
	"github.com/bdv/gprs/pkg/config"
	"github.com/bdv/gprs/pkg/id"
	"google.golang.org/grpc"
	healthgrpc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.LoadPayments()

	repo := paymentsrepository.NewInMemoryRepository()
	idGen := id.NewGenerator()
	service := paymentsservice.NewService(repo, idGen)

	grpcServer := grpc.NewServer()
	paymentspb.RegisterPaymentsServiceServer(grpcServer, paymentsgrpc.NewGRPCServer(service))

	healthServer := healthgrpc.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", cfg.GRPCPort, err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("received signal %s, shutting down payments service", sig)
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		grpcServer.GracefulStop()
	}()

	log.Printf("payments service listening on %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("payments service stopped: %v", err)
	}
}
