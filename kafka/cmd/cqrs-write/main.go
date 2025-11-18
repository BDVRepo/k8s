package main

import (
	"context"
	"log"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/patterns/cqrs"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	producer := kafka.NewProducer(cfg.Brokers)
	defer producer.Close()

	writeSide := cqrs.NewWriteSide(producer)
	ctx := context.Background()

	log.Println("CQRS Write Side started...")

	// Цикл демонстрирует сценарий:
	// 1) команда CreateOrder -> публикуем OrderCreated
	// 2) команда ProcessPayment -> публикуем OrderPaid
	// Read-side ничего об этом не знает, он просто слушает Kafka.
	for i := 0; i < 5; i++ {
		orderID, err := writeSide.CreateOrder(ctx, cqrs.CreateOrderCommand{
			UserID:    "user-123",
			ProductID: "product-456",
			Amount:    99.99 + float64(i)*10,
			Currency:  "USD",
		})
		if err != nil {
			log.Printf("Failed to create order: %v", err)
			continue
		}

		log.Printf("Created order: %s", orderID)

		// Имитируем задержку между командами (как разные сервисы).
		time.Sleep(1 * time.Second)
		paymentID, err := writeSide.ProcessPayment(ctx, cqrs.ProcessPaymentCommand{
			OrderID: orderID,
			Amount:  99.99 + float64(i)*10,
		})
		if err != nil {
			log.Printf("Failed to process payment: %v", err)
			continue
		}

		log.Printf("Processed payment: %s for order: %s", paymentID, orderID)
		time.Sleep(1 * time.Second)
	}

	log.Println("CQRS Write Side example completed")
}
