package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
	"github.com/bdv/kafka-learning/internal/patterns/cqrs"
	"github.com/bdv/kafka-learning/pkg/config"
)

func main() {
	cfg := config.LoadKafka()
	// Read side слушает два стрима независимо.
	ordersConsumer := kafka.NewConsumer(cfg.Brokers, "orders", "analytics-consumer-group")
	defer ordersConsumer.Close()

	paymentsConsumer := kafka.NewConsumer(cfg.Brokers, "payments", "analytics-consumer-group")
	defer paymentsConsumer.Close()

	// ReadModel хранит агрегаты, которые невозможно быстро получить из OLTP.
	readModel := cqrs.NewAnalyticsReadModel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigCh
		log.Printf("Received signal %s, shutting down analytics consumer", sig)
		cancel()
	}()

	// Периодически показываем статистику
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				totalOrders, totalRevenue, users, products := readModel.GetStats()
				log.Printf("[analytics] Stats - Total Orders: %d, Total Revenue: %.2f",
					totalOrders, totalRevenue)
				log.Printf("[analytics] Top Users: %v", users)
				log.Printf("[analytics] Top Products: %v", products)

				recent := readModel.GetRecentOrders(5)
				log.Printf("[analytics] Recent Orders: %d", len(recent))
				for _, order := range recent {
					log.Printf("  - %s: %s, %.2f, %s", order.OrderID, order.UserID, order.Amount, order.Status)
				}

			case <-ctx.Done():
				return
			}
		}
	}()

	log.Println("Analytics Consumer (CQRS Read Side) started, waiting for messages...")

	// Поток #1: события создания заказов → обновляем счетчики, revenue и recent orders.
	go func() {
		for {
			msg, err := ordersConsumer.ReadMessage(ctx)
			if err != nil {
				if err == context.Canceled {
					return
				}
				log.Printf("Error reading orders message: %v", err)
				continue
			}

			ordersConsumer.LogMessage(msg)

			var event models.OrderCreatedEvent
			if err := json.Unmarshal(msg.Value, &event); err != nil {
				log.Printf("Error unmarshaling OrderCreated event: %v", err)
				continue
			}

			readModel.HandleOrderCreated(ctx, event)
		}
	}()

	// Поток #2: события оплаты → меняем статус заказов в read model.
	for {
		msg, err := paymentsConsumer.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Println("Analytics consumer stopped")
				return
			}
			log.Printf("Error reading payments message: %v", err)
			continue
		}

		paymentsConsumer.LogMessage(msg)

		var event models.OrderPaidEvent
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Error unmarshaling OrderPaid event: %v", err)
			continue
		}

		readModel.HandleOrderPaid(ctx, event)
	}
}
