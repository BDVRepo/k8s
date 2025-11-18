package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/patterns/outbox"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	producer := kafka.NewProducer(cfg.Brokers)
	defer producer.Close()

	// OutboxStore имитирует таблицу outbox (map[eventID]OutboxEvent).
	outboxStore := outbox.NewInMemoryOutboxStore()
	// TransactionalService делает "Insert order + Insert outbox" под одной транзакцией.
	service := outbox.NewTransactionalService(outboxStore)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Publisher — отдельный процесс, который опрашивает outbox и публикует в Kafka.
	publisher := outbox.NewOutboxPublisher(outboxStore, producer, 2*time.Second)
	go publisher.Start(ctx)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal %s, shutting down", sig)
		cancel()
	}()

	log.Println("Outbox Pattern example started...")
	log.Println("Creating orders with outbox pattern...")

	// Application-часть: создаем заказы, и каждый вызов CreateOrderWithOutbox
	// сохраняет событие в outbox (но не отправляет напрямую в Kafka).
	for i := 0; i < 5; i++ {
		orderID := "order-outbox-" + time.Now().Format("150405") + string(rune('0'+i))
		if err := service.CreateOrderWithOutbox(ctx, orderID, "user-123", "product-456", 99.99+float64(i)*10); err != nil {
			log.Printf("Failed to create order: %v", err)
			continue
		}

		time.Sleep(1 * time.Second)
	}

	log.Println("Orders created and saved to outbox. Publisher will publish them to Kafka once available...")
	log.Println("Press Ctrl+C to stop")

	// Ждем завершения
	<-ctx.Done()
	log.Println("Outbox Pattern example completed")
}
