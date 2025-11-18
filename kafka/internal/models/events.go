package models

import "time"

// OrderCreatedEvent событие создания заказа
type OrderCreatedEvent struct {
	OrderID   string    `json:"order_id"`
	UserID    string    `json:"user_id"`
	ProductID string    `json:"product_id"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

// OrderPaidEvent событие оплаты заказа
type OrderPaidEvent struct {
	OrderID   string    `json:"order_id"`
	PaymentID string    `json:"payment_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

// OrderShippedEvent событие отправки заказа
type OrderShippedEvent struct {
	OrderID   string    `json:"order_id"`
	TrackingNumber string `json:"tracking_number"`
	Timestamp time.Time `json:"timestamp"`
}

