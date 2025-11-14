package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	orderspb "github.com/bdv/gprs/internal/gen/orderspb"
	paymentspb "github.com/bdv/gprs/internal/gen/paymentspb"
	grpcclient "github.com/bdv/gprs/internal/infrastructure/grpc"
	ordersrepository "github.com/bdv/gprs/internal/methods/handlers/orders/adapter"
	ordersgrpc "github.com/bdv/gprs/internal/methods/handlers/orders/grpc"
	ordersservice "github.com/bdv/gprs/internal/methods/handlers/orders/usecase"
	"github.com/bdv/gprs/pkg/config"
	"github.com/bdv/gprs/pkg/id"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthgrpc "google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadOrders()

	conn, err := grpc.DialContext(ctx, cfg.PaymentsAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to payments service: %v", err)
	}
	defer conn.Close()

	paymentClient := grpcclient.NewGRPCPaymentClient(paymentspb.NewPaymentsServiceClient(conn))
	repo := ordersrepository.NewInMemoryRepository()
	idGen := id.NewGenerator()
	service := ordersservice.NewService(repo, paymentClient, idGen)

	grpcServer := grpc.NewServer()
	orderspb.RegisterOrdersServiceServer(grpcServer, ordersgrpc.NewGRPCServer(service))

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
		log.Printf("received signal %s, shutting down orders service", sig)
		healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
		grpcServer.GracefulStop()
	}()

	log.Printf("orders service listening on %s", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("orders service stopped: %v", err)
	}
}
