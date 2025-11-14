package config

import "os"

type OrdersConfig struct {
	GRPCPort     string
	PaymentsAddr string
}

type PaymentsConfig struct {
	GRPCPort string
}

func LoadOrders() OrdersConfig {
	return OrdersConfig{
		GRPCPort:     getEnv("ORDERS_GRPC_PORT", ":50051"),
		PaymentsAddr: getEnv("PAYMENTS_GRPC_ADDR", "localhost:50052"),
	}
}

func LoadPayments() PaymentsConfig {
	return PaymentsConfig{
		GRPCPort: getEnv("PAYMENTS_GRPC_PORT", ":50052"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
