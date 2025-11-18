package main

import (
	"context"
	"log"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
	"github.com/bdv/kafka-learning/internal/patterns/eventsourcing"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	producer := kafka.NewProducer(cfg.Brokers)
	defer producer.Close()

	eventStore := eventsourcing.NewInMemoryEventStore(producer)
	ctx := context.Background()

	log.Println("Orders Producer (Event Sourcing) started...")

	// 1) Создаем заказ и генерируем первое событие OrderCreated.
	//    Event Store добавит его в map[orderID][]Event и отправит в Kafka (topic orders).
	orderID := "order-es-001"
	createdEvent := models.OrderCreatedEvent{
		OrderID:   orderID,
		UserID:    "user-123",
		ProductID: "product-456",
		Amount:    199.99,
		Currency:  "USD",
		Timestamp: time.Now(),
	}

	if err := eventStore.SaveEvent(ctx, createdEvent); err != nil {
		log.Fatalf("Failed to save OrderCreated event: %v", err)
	}

	// Каждый RebuildState проходит по всем событиям orderID и "переигрывает" их,
	// чтобы собрать актуальное состояние (без хранения отдельной записи в БД).
	state, err := eventStore.RebuildState(ctx, orderID)
	if err != nil {
		log.Fatalf("Failed to rebuild state: %v", err)
	}
	log.Printf("State after creation: OrderID=%s, Status=%s", state.OrderID, state.Status)

	time.Sleep(2 * time.Second)

	// 2) Создаем событие OrderPaid – снова кладем в Event Store.
	paidEvent := models.OrderPaidEvent{
		OrderID:   orderID,
		PaymentID: "payment-789",
		Amount:    199.99,
		Timestamp: time.Now(),
	}

	if err := eventStore.SaveEvent(ctx, paidEvent); err != nil {
		log.Fatalf("Failed to save OrderPaid event: %v", err)
	}

	// После второго события RebuildState видит Created + Paid и меняет статус на PAID.
	state, err = eventStore.RebuildState(ctx, orderID)
	if err != nil {
		log.Fatalf("Failed to rebuild state: %v", err)
	}
	log.Printf("State after payment: OrderID=%s, Status=%s, PaymentID=%s",
		state.OrderID, state.Status, state.PaymentID)

	time.Sleep(2 * time.Second)

	// 3) Завершаем сценарий событием OrderShipped (topic shipping).
	shippedEvent := models.OrderShippedEvent{
		OrderID:        orderID,
		TrackingNumber: "TRACK-123456",
		Timestamp:      time.Now(),
	}

	if err := eventStore.SaveEvent(ctx, shippedEvent); err != nil {
		log.Fatalf("Failed to save OrderShipped event: %v", err)
	}

	// Финальный rebuild покажет статус SHIPPED и трекинг.
	state, err = eventStore.RebuildState(ctx, orderID)
	if err != nil {
		log.Fatalf("Failed to rebuild state: %v", err)
	}
	log.Printf("Final state: OrderID=%s, Status=%s, TrackingNumber=%s",
		state.OrderID, state.Status, state.TrackingNumber)

	// Показываем все события
	events, err := eventStore.GetEvents(ctx, orderID)
	if err != nil {
		log.Fatalf("Failed to get events: %v", err)
	}
	log.Printf("Total events for order %s: %d", orderID, len(events))

	log.Println("Event Sourcing example completed")
}
