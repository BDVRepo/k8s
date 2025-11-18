package models

import "time"

// SagaStepCompletedEvent событие завершения шага Saga
type SagaStepCompletedEvent struct {
	SagaID    string    `json:"saga_id"`
	StepName  string    `json:"step_name"`
	OrderID   string    `json:"order_id"`
	Success   bool      `json:"success"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// SagaCompensationEvent событие компенсации шага Saga
type SagaCompensationEvent struct {
	SagaID    string    `json:"saga_id"`
	StepName  string    `json:"step_name"`
	OrderID   string    `json:"order_id"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

// SagaCompletedEvent событие завершения всей Saga
type SagaCompletedEvent struct {
	SagaID    string    `json:"saga_id"`
	OrderID   string    `json:"order_id"`
	Success   bool      `json:"success"`
	Timestamp time.Time `json:"timestamp"`
}

func (e SagaStepCompletedEvent) GetOrderID() string {
	return e.OrderID
}

func (e SagaStepCompletedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e SagaCompensationEvent) GetOrderID() string {
	return e.OrderID
}

func (e SagaCompensationEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e SagaCompletedEvent) GetOrderID() string {
	return e.OrderID
}

func (e SagaCompletedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

