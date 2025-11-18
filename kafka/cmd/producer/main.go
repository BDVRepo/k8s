package main

import (
	"context"
	"log"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	producer := kafka.NewProducer(cfg.Brokers)
	defer producer.Close()

	ctx := context.Background()
	topic := "orders"

	log.Println("Producer started, publishing messages...")

	for i := 0; i < 10; i++ {
		event := models.OrderCreatedEvent{
			OrderID:   "order-" + string(rune('0'+i)),
			UserID:    "user-123",
			ProductID: "product-456",
			Amount:    99.99 + float64(i),
			Currency:  "USD",
			Timestamp: time.Now(),
		}

		key := event.OrderID
		if err := producer.Publish(ctx, topic, key, event); err != nil {
			log.Printf("Failed to publish event: %v", err)
			continue
		}

		log.Printf("Published event: OrderID=%s", event.OrderID)
		time.Sleep(1 * time.Second)
	}

	log.Println("Producer finished")
}
