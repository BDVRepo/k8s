package eventsourcing

import (
	"context"
	"log"
	"sync"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
)

// EventStore хранит все события для Event Sourcing
type EventStore interface {
	SaveEvent(ctx context.Context, event models.Event) error
	GetEvents(ctx context.Context, orderID string) ([]models.Event, error)
	RebuildState(ctx context.Context, orderID string) (*models.OrderState, error)
}

// InMemoryEventStore in-memory реализация для примера
type InMemoryEventStore struct {
	mu       sync.RWMutex
	events   map[string][]models.Event // orderID -> events
	producer *kafka.Producer
}

func NewInMemoryEventStore(producer *kafka.Producer) *InMemoryEventStore {
	return &InMemoryEventStore{
		events:   make(map[string][]models.Event),
		producer: producer,
	}
}

func (s *InMemoryEventStore) SaveEvent(ctx context.Context, event models.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	orderID := event.GetOrderID()
	s.events[orderID] = append(s.events[orderID], event)

	// Определяем топик по типу события
	topic := s.getTopicForEvent(event)
	if err := s.producer.Publish(ctx, topic, orderID, event); err != nil {
		log.Printf("[event-store] failed to publish event to kafka: %v", err)
		// Не возвращаем ошибку, т.к. событие уже сохранено локально
	}

	log.Printf("[event-store] saved event for order_id=%s, total_events=%d",
		orderID, len(s.events[orderID]))

	return nil
}

func (s *InMemoryEventStore) GetEvents(ctx context.Context, orderID string) ([]models.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events, ok := s.events[orderID]
	if !ok {
		return nil, nil
	}

	return events, nil
}

func (s *InMemoryEventStore) RebuildState(ctx context.Context, orderID string) (*models.OrderState, error) {
	events, err := s.GetEvents(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, nil
	}

	state := &models.OrderState{}

	for _, event := range events {
		switch e := event.(type) {
		case models.OrderCreatedEvent:
			state.OrderID = e.OrderID
			state.UserID = e.UserID
			state.ProductID = e.ProductID
			state.Amount = e.Amount
			state.Currency = e.Currency
			state.Status = "CREATED"
			state.CreatedAt = e.Timestamp

		case models.OrderPaidEvent:
			state.PaymentID = e.PaymentID
			state.Status = "PAID"
			paidAt := e.Timestamp
			state.PaidAt = &paidAt

		case models.OrderShippedEvent:
			state.TrackingNumber = e.TrackingNumber
			state.Status = "SHIPPED"
			shippedAt := e.Timestamp
			state.ShippedAt = &shippedAt
		}
	}

	log.Printf("[event-store] rebuilt state for order_id=%s, status=%s, events_count=%d",
		orderID, state.Status, len(events))

	return state, nil
}

func (s *InMemoryEventStore) getTopicForEvent(event models.Event) string {
	switch event.(type) {
	case models.OrderCreatedEvent:
		return "orders"
	case models.OrderPaidEvent:
		return "payments"
	case models.OrderShippedEvent:
		return "shipping"
	default:
		return "events"
	}
}
