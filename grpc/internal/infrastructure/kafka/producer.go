package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type EventPublisher interface {
	PublishOrderCreated(ctx context.Context, orderID, userID, productID string, amount float64, currency string) error
}

type KafkaEventPublisher struct {
	writer *kafka.Writer
}

func NewKafkaEventPublisher(brokers []string) *KafkaEventPublisher {
	return &KafkaEventPublisher{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

type OrderCreatedEvent struct {
	OrderID   string  `json:"order_id"`
	UserID    string  `json:"user_id"`
	ProductID string  `json:"product_id"`
	Amount    float64 `json:"amount"`
	Currency  string  `json:"currency"`
}

func (p *KafkaEventPublisher) PublishOrderCreated(ctx context.Context, orderID, userID, productID string, amount float64, currency string) error {
	event := OrderCreatedEvent{
		OrderID:   orderID,
		UserID:    userID,
		ProductID: productID,
		Amount:    amount,
		Currency:  currency,
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(orderID),
		Value: eventJSON,
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("[kafka-publisher] failed to publish OrderCreated event: %v", err)
		return err
	}

	log.Printf("[kafka-publisher] published OrderCreated event: order_id=%s", orderID)
	return nil
}

func (p *KafkaEventPublisher) Close() error {
	return p.writer.Close()
}

