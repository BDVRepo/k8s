package models

import "time"

// OrderState состояние заказа, восстанавливается из событий
type OrderState struct {
	OrderID       string
	UserID        string
	ProductID     string
	Amount        float64
	Currency      string
	Status        string
	PaymentID     string
	TrackingNumber string
	CreatedAt     time.Time
	PaidAt        *time.Time
	ShippedAt     *time.Time
}

// Event интерфейс для всех событий
type Event interface {
	GetOrderID() string
	GetTimestamp() time.Time
}

func (e OrderCreatedEvent) GetOrderID() string {
	return e.OrderID
}

func (e OrderCreatedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e OrderPaidEvent) GetOrderID() string {
	return e.OrderID
}

func (e OrderPaidEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

func (e OrderShippedEvent) GetOrderID() string {
	return e.OrderID
}

func (e OrderShippedEvent) GetTimestamp() time.Time {
	return e.Timestamp
}

