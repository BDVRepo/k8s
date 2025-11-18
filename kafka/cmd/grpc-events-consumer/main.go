package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/pkg/config"
)

type OrderCreatedEvent struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

func main() {
	cfg := config.LoadKafka()
	consumer := kafka.NewConsumer(cfg.Brokers, "orders", "grpc-events-consumer-group")
	defer consumer.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal %s, shutting down consumer", sig)
		cancel()
	}()

	log.Println("gRPC Events Consumer started, listening for OrderCreated events from gRPC service...")

	for {
		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Println("Consumer stopped")
				return
			}
			log.Printf("Error reading message: %v", err)
			continue
		}

		consumer.LogMessage(msg)

		var event OrderCreatedEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling event: %v", err)
			continue
		}

		log.Printf("[grpc-events-consumer] received OrderCreated from gRPC service:")
		log.Printf("  OrderID: %s", event.OrderID)
		log.Printf("  UserID: %s", event.UserID)
		log.Printf("  ProductID: %s", event.ProductID)
		log.Printf("  Amount: %.2f %s", event.Amount, event.Currency)

		// Здесь можно добавить обработку события:
		// - Отправка email
		// - Обновление аналитики
		// - Индексация для поиска
		// и т.д.
	}
}

