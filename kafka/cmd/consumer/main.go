package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	consumer := kafka.NewConsumer(cfg.Brokers, "orders", "orders-consumer-group")
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

	log.Println("Consumer started, waiting for messages...")

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

		var event models.OrderCreatedEvent
		if err := consumer.UnmarshalMessage(msg, &event); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		log.Printf("Processed event: OrderID=%s, UserID=%s, Amount=%.2f %s",
			event.OrderID, event.UserID, event.Amount, event.Currency)

		// Здесь можно добавить бизнес-логику обработки события
	}
}
