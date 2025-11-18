package cqrs

import (
	"context"
	"log"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
)

// WriteSide обрабатывает команды и публикует события
type WriteSide struct {
	producer *kafka.Producer
}

func NewWriteSide(producer *kafka.Producer) *WriteSide {
	return &WriteSide{
		producer: producer,
	}
}

// CreateOrderCommand команда для создания заказа
type CreateOrderCommand struct {
	UserID    string
	ProductID string
	Amount    float64
	Currency  string
}

// CreateOrder обрабатывает команду создания заказа
func (w *WriteSide) CreateOrder(ctx context.Context, cmd CreateOrderCommand) (string, error) {
	// Генерируем ID заказа
	orderID := "order-" + generateID()

	// Создаем событие
	event := models.OrderCreatedEvent{
		OrderID:   orderID,
		UserID:    cmd.UserID,
		ProductID: cmd.ProductID,
		Amount:    cmd.Amount,
		Currency:  cmd.Currency,
		Timestamp: getCurrentTime(),
	}

	// Публикуем событие в Kafka (write side только публикует события)
	if err := w.producer.Publish(ctx, "orders", orderID, event); err != nil {
		return "", err
	}

	log.Printf("[write-side] published OrderCreated event: order_id=%s", orderID)
	return orderID, nil
}

// ProcessPaymentCommand команда для обработки платежа
type ProcessPaymentCommand struct {
	OrderID string
	Amount  float64
}

// ProcessPayment обрабатывает команду платежа
func (w *WriteSide) ProcessPayment(ctx context.Context, cmd ProcessPaymentCommand) (string, error) {
	paymentID := "payment-" + generateID()

	event := models.OrderPaidEvent{
		OrderID:   cmd.OrderID,
		PaymentID: paymentID,
		Amount:    cmd.Amount,
		Timestamp: getCurrentTime(),
	}

	if err := w.producer.Publish(ctx, "payments", cmd.OrderID, event); err != nil {
		return "", err
	}

	log.Printf("[write-side] published OrderPaid event: order_id=%s, payment_id=%s", cmd.OrderID, paymentID)
	return paymentID, nil
}

