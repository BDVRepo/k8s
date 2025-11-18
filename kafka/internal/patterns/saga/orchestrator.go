package saga

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/bdv/kafka-learning/internal/kafka"
	"github.com/bdv/kafka-learning/internal/models"
)

// SagaOrchestrator оркестрирует распределенную транзакцию
type SagaOrchestrator struct {
	producer *kafka.Producer
	sagaID   string
	orderID  string
	steps    []SagaStep
}

// SagaStep шаг в Saga
type SagaStep struct {
	Name         string
	ActionTopic  string
	CompensateTopic string
	Action       func(ctx context.Context) error
	Compensate   func(ctx context.Context) error
}

// NewSagaOrchestrator создает новый orchestrator
func NewSagaOrchestrator(producer *kafka.Producer, sagaID, orderID string) *SagaOrchestrator {
	return &SagaOrchestrator{
		producer: producer,
		sagaID:   sagaID,
		orderID:  orderID,
		steps:    make([]SagaStep, 0),
	}
}

// AddStep добавляет шаг в Saga
func (s *SagaOrchestrator) AddStep(name, actionTopic, compensateTopic string) {
	s.steps = append(s.steps, SagaStep{
		Name:          name,
		ActionTopic:   actionTopic,
		CompensateTopic: compensateTopic,
	})
}

// Execute выполняет Saga с компенсацией при ошибках
func (s *SagaOrchestrator) Execute(ctx context.Context) error {
	log.Printf("[saga] starting saga_id=%s order_id=%s", s.sagaID, s.orderID)

	var completedSteps []int

	for i, step := range s.steps {
		log.Printf("[saga] executing step %d: %s", i+1, step.Name)

		// Публикуем команду для выполнения шага
		command := map[string]interface{}{
			"saga_id":  s.sagaID,
			"order_id": s.orderID,
			"step":     step.Name,
		}

		commandJSON, _ := json.Marshal(command)
		if err := s.producer.Publish(ctx, step.ActionTopic, s.orderID, commandJSON); err != nil {
			log.Printf("[saga] failed to publish command for step %s: %v", step.Name, err)
			// Компенсируем предыдущие шаги
			s.compensate(ctx, completedSteps)
			return err
		}

		// В реальном приложении здесь нужно ждать события SagaStepCompletedEvent
		// Для примера просто ждем немного
		time.Sleep(500 * time.Millisecond)

		// Публикуем событие завершения шага (в реальности это делает сервис)
		stepEvent := models.SagaStepCompletedEvent{
			SagaID:    s.sagaID,
			StepName:  step.Name,
			OrderID:   s.orderID,
			Success:   true,
			Timestamp: time.Now(),
		}

		if err := s.producer.Publish(ctx, "saga-events", s.sagaID, stepEvent); err != nil {
			log.Printf("[saga] failed to publish step completion event: %v", err)
		}

		completedSteps = append(completedSteps, i)
		log.Printf("[saga] step %s completed", step.Name)
	}

	// Публикуем событие завершения Saga
	completedEvent := models.SagaCompletedEvent{
		SagaID:    s.sagaID,
		OrderID:   s.orderID,
		Success:   true,
		Timestamp: time.Now(),
	}

	if err := s.producer.Publish(ctx, "saga-events", s.sagaID, completedEvent); err != nil {
		log.Printf("[saga] failed to publish saga completion event: %v", err)
	}

	log.Printf("[saga] saga completed successfully: saga_id=%s", s.sagaID)
	return nil
}

// ExecuteWithFailure выполняет Saga с симуляцией ошибки
func (s *SagaOrchestrator) ExecuteWithFailure(ctx context.Context, failAtStep int) error {
	log.Printf("[saga] starting saga with failure simulation: saga_id=%s, fail_at_step=%d", s.sagaID, failAtStep)

	var completedSteps []int

	for i, step := range s.steps {
		if i == failAtStep {
			log.Printf("[saga] simulating failure at step %d: %s", i+1, step.Name)

			// Публикуем событие ошибки
			stepEvent := models.SagaStepCompletedEvent{
				SagaID:    s.sagaID,
				StepName:  step.Name,
				OrderID:   s.orderID,
				Success:   false,
				Error:     "simulated error",
				Timestamp: time.Now(),
			}

			s.producer.Publish(ctx, "saga-events", s.sagaID, stepEvent)

			// Компенсируем предыдущие шаги
			s.compensate(ctx, completedSteps)
			return nil
		}

		log.Printf("[saga] executing step %d: %s", i+1, step.Name)

		command := map[string]interface{}{
			"saga_id":  s.sagaID,
			"order_id": s.orderID,
			"step":     step.Name,
		}

		commandJSON, _ := json.Marshal(command)
		s.producer.Publish(ctx, step.ActionTopic, s.orderID, commandJSON)

		time.Sleep(500 * time.Millisecond)

		stepEvent := models.SagaStepCompletedEvent{
			SagaID:    s.sagaID,
			StepName:  step.Name,
			OrderID:   s.orderID,
			Success:   true,
			Timestamp: time.Now(),
		}

		s.producer.Publish(ctx, "saga-events", s.sagaID, stepEvent)
		completedSteps = append(completedSteps, i)
	}

	return nil
}

// compensate выполняет компенсацию для завершенных шагов
func (s *SagaOrchestrator) compensate(ctx context.Context, completedSteps []int) {
	log.Printf("[saga] starting compensation for %d steps", len(completedSteps))

	// Компенсируем в обратном порядке
	for i := len(completedSteps) - 1; i >= 0; i-- {
		stepIndex := completedSteps[i]
		step := s.steps[stepIndex]

		log.Printf("[saga] compensating step: %s", step.Name)

		compensationEvent := models.SagaCompensationEvent{
			SagaID:    s.sagaID,
			StepName:  step.Name,
			OrderID:   s.orderID,
			Reason:    "saga failed, rolling back",
			Timestamp: time.Now(),
		}

		s.producer.Publish(ctx, step.CompensateTopic, s.orderID, compensationEvent)
		time.Sleep(200 * time.Millisecond)
	}

	log.Printf("[saga] compensation completed")
}

