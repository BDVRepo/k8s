.PHONY: proto tidy build-orders build-payments run-orders run-payments test

PROTOC_GEN := $(shell go env GOPATH)/bin
PROTO_SRC := $(shell find proto -name '*.proto')

proto: ## Generate gRPC stubs from protobuf definitions
	PATH="$(PATH):$(PROTOC_GEN)" protoc --proto_path=proto --go_out=internal/gen --go_opt=paths=source_relative --go-grpc_out=internal/gen --go-grpc_opt=paths=source_relative $(PROTO_SRC)
	@rm -rf internal/gen/orderspb internal/gen/paymentspb
	@mkdir -p internal/gen/orderspb internal/gen/paymentspb
	@mv internal/gen/orders/*.go internal/gen/orderspb/
	@mv internal/gen/payments/*.go internal/gen/paymentspb/
	@rmdir internal/gen/orders internal/gen/payments

 tidy: ## Ensure Go module deps are tidy
	go mod tidy

build-orders: ## Build the orders service binary
	go build -o bin/orders ./cmd/orders

build-payments: ## Build the payments service binary
	go build -o bin/payments ./cmd/payments

run-payments: ## Run the payments service locally
	PAYMENTS_GRPC_PORT=:50052 go run ./cmd/payments

run-orders: ## Run the orders service locally
	ORDERS_GRPC_PORT=:50051 PAYMENTS_GRPC_ADDR=localhost:50052 go run ./cmd/orders

test: ## Execute unit tests
	go test ./...

help: ## Display available targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'
