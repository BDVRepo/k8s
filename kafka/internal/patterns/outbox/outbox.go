package outbox

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
)

// OutboxEvent событие в outbox таблице
type OutboxEvent struct {
	ID        string
	Topic     string
	Key       string
	Value     []byte
	CreatedAt time.Time
	Published bool
}

// OutboxStore хранит события в outbox
type OutboxStore interface {
	SaveEvent(ctx context.Context, event OutboxEvent) error
	GetUnpublishedEvents(ctx context.Context, limit int) ([]OutboxEvent, error)
	MarkAsPublished(ctx context.Context, eventID string) error
}

// InMemoryOutboxStore in-memory реализация для примера
type InMemoryOutboxStore struct {
	mu     sync.RWMutex
	events map[string]OutboxEvent
}

func NewInMemoryOutboxStore() *InMemoryOutboxStore {
	return &InMemoryOutboxStore{
		events: make(map[string]OutboxEvent),
	}
}

func (s *InMemoryOutboxStore) SaveEvent(ctx context.Context, event OutboxEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	log.Printf("[outbox] saved event to outbox: id=%s, topic=%s, key=%s", event.ID, event.Topic, event.Key)
	return nil
}

func (s *InMemoryOutboxStore) GetUnpublishedEvents(ctx context.Context, limit int) ([]OutboxEvent, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var unpublished []OutboxEvent
	for _, event := range s.events {
		if !event.Published {
			unpublished = append(unpublished, event)
			if len(unpublished) >= limit {
				break
			}
		}
	}

	return unpublished, nil
}

func (s *InMemoryOutboxStore) MarkAsPublished(ctx context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, ok := s.events[eventID]
	if !ok {
		return nil
	}

	event.Published = true
	s.events[eventID] = event
	log.Printf("[outbox] marked as published: id=%s", eventID)
	return nil
}

// OutboxPublisher публикует события из outbox в Kafka
type OutboxPublisher struct {
	store    OutboxStore
	producer *kafka.Producer
	interval time.Duration
}

func NewOutboxPublisher(store OutboxStore, producer *kafka.Producer, interval time.Duration) *OutboxPublisher {
	return &OutboxPublisher{
		store:    store,
		producer: producer,
		interval: interval,
	}
}

// Start запускает publisher в фоне
func (p *OutboxPublisher) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	log.Println("[outbox-publisher] started")

	for {
		select {
		case <-ticker.C:
			p.publishPendingEvents(ctx)
		case <-ctx.Done():
			log.Println("[outbox-publisher] stopped")
			return
		}
	}
}

func (p *OutboxPublisher) publishPendingEvents(ctx context.Context) {
	events, err := p.store.GetUnpublishedEvents(ctx, 10)
	if err != nil {
		log.Printf("[outbox-publisher] failed to get unpublished events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	log.Printf("[outbox-publisher] found %d unpublished events", len(events))

	for _, event := range events {
		// Публикуем в Kafka
		if err := p.producer.Publish(ctx, event.Topic, event.Key, event.Value); err != nil {
			log.Printf("[outbox-publisher] failed to publish event %s: %v", event.ID, err)
			continue
		}

		// Помечаем как опубликованное
		if err := p.store.MarkAsPublished(ctx, event.ID); err != nil {
			log.Printf("[outbox-publisher] failed to mark event %s as published: %v", event.ID, err)
			continue
		}

		log.Printf("[outbox-publisher] published event: id=%s, topic=%s", event.ID, event.Topic)
	}
}

// TransactionalService сервис с поддержкой транзакций и outbox
type TransactionalService struct {
	outbox OutboxStore
}

func NewTransactionalService(outbox OutboxStore) *TransactionalService {
	return &TransactionalService{
		outbox: outbox,
	}
}

// CreateOrderWithOutbox создает заказ и сохраняет событие в outbox в одной транзакции
func (s *TransactionalService) CreateOrderWithOutbox(ctx context.Context, orderID, userID, productID string, amount float64) error {
	// В реальном приложении здесь была бы БД транзакция
	// BEGIN TRANSACTION
	//   INSERT INTO orders ...
	//   INSERT INTO outbox ...
	// COMMIT

	// Для примера просто сохраняем в outbox
	eventData := map[string]interface{}{
		"order_id":   orderID,
		"user_id":    userID,
		"product_id": productID,
		"amount":     amount,
		"timestamp":  time.Now(),
	}

	eventJSON, _ := json.Marshal(eventData)

	outboxEvent := OutboxEvent{
		ID:        "outbox-" + orderID,
		Topic:     "orders",
		Key:       orderID,
		Value:     eventJSON,
		CreatedAt: time.Now(),
		Published: false,
	}

	if err := s.outbox.SaveEvent(ctx, outboxEvent); err != nil {
		return err
	}

	log.Printf("[transactional-service] created order %s and saved to outbox", orderID)
	return nil
}

