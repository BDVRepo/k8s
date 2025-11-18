package main

import (
	"context"
	"log"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/patterns/saga"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	producer := kafka.NewProducer(cfg.Brokers)
	defer producer.Close()

	ctx := context.Background()

	// Пример 1: успешная saga (оркестратор выполняет три шага без ошибок).
	// Каждый AddStep принимает: название шага, topic для команды, topic для компенсации.
	log.Println("=== Example 1: Successful Saga ===")
	sagaID1 := "saga-" + time.Now().Format("20060102150405")
	orderID1 := "order-saga-001"

	orchestrator1 := saga.NewSagaOrchestrator(producer, sagaID1, orderID1)
	orchestrator1.AddStep("CreateOrder", "order-commands", "order-compensate")
	orchestrator1.AddStep("ProcessPayment", "payment-commands", "payment-compensate")
	orchestrator1.AddStep("ShipOrder", "shipping-commands", "shipping-compensate")

	if err := orchestrator1.Execute(ctx); err != nil {
		log.Printf("Saga failed: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Пример 2: симулируем ошибку на шаге ProcessPayment (index 1),
	// оркестратор вызывает компенсации для уже выполненных шагов.
	log.Println("\n=== Example 2: Saga with Failure and Compensation ===")
	sagaID2 := "saga-" + time.Now().Format("20060102150405")
	orderID2 := "order-saga-002"

	orchestrator2 := saga.NewSagaOrchestrator(producer, sagaID2, orderID2)
	orchestrator2.AddStep("CreateOrder", "order-commands", "order-compensate")
	orchestrator2.AddStep("ProcessPayment", "payment-commands", "payment-compensate")
	orchestrator2.AddStep("ShipOrder", "shipping-commands", "shipping-compensate")

	// Симулируем ошибку на втором шаге (ProcessPayment)
	orchestrator2.ExecuteWithFailure(ctx, 1)

	log.Println("\nSaga Pattern examples completed")
}
