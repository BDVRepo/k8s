package cqrs

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/bdv/kafka-learning/internal/models"
)

// AnalyticsReadModel read model для аналитики (оптимизирован для чтения)
type AnalyticsReadModel struct {
	mu              sync.RWMutex
	totalOrders     int
	totalRevenue    float64
	ordersByUser    map[string]int
	ordersByProduct map[string]int
	recentOrders    []OrderSummary
}

type OrderSummary struct {
	OrderID   string
	UserID    string
	ProductID string
	Amount    float64
	Status    string
	Timestamp time.Time
}

func NewAnalyticsReadModel() *AnalyticsReadModel {
	return &AnalyticsReadModel{
		ordersByUser:    make(map[string]int),
		ordersByProduct: make(map[string]int),
		recentOrders:    make([]OrderSummary, 0),
	}
}

// HandleOrderCreated обрабатывает событие создания заказа
func (r *AnalyticsReadModel) HandleOrderCreated(ctx context.Context, event models.OrderCreatedEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.totalOrders++
	r.totalRevenue += event.Amount
	r.ordersByUser[event.UserID]++
	r.ordersByProduct[event.ProductID]++

	// Добавляем в список недавних заказов
	summary := OrderSummary{
		OrderID:   event.OrderID,
		UserID:    event.UserID,
		ProductID: event.ProductID,
		Amount:    event.Amount,
		Status:    "CREATED",
		Timestamp: event.Timestamp,
	}

	r.recentOrders = append(r.recentOrders, summary)
	if len(r.recentOrders) > 100 {
		r.recentOrders = r.recentOrders[1:] // Оставляем последние 100
	}

	log.Printf("[read-side] updated analytics: total_orders=%d, total_revenue=%.2f",
		r.totalOrders, r.totalRevenue)
}

// HandleOrderPaid обрабатывает событие оплаты заказа
func (r *AnalyticsReadModel) HandleOrderPaid(ctx context.Context, event models.OrderPaidEvent) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Обновляем статус в recentOrders
	for i := range r.recentOrders {
		if r.recentOrders[i].OrderID == event.OrderID {
			r.recentOrders[i].Status = "PAID"
			break
		}
	}

	log.Printf("[read-side] updated order status: order_id=%s, status=PAID", event.OrderID)
}

// GetStats возвращает статистику (read операция)
func (r *AnalyticsReadModel) GetStats() (int, float64, map[string]int, map[string]int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Копируем maps для безопасного возврата
	usersCopy := make(map[string]int)
	productsCopy := make(map[string]int)

	for k, v := range r.ordersByUser {
		usersCopy[k] = v
	}
	for k, v := range r.ordersByProduct {
		productsCopy[k] = v
	}

	return r.totalOrders, r.totalRevenue, usersCopy, productsCopy
}

// GetRecentOrders возвращает недавние заказы
func (r *AnalyticsReadModel) GetRecentOrders(limit int) []OrderSummary {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if limit > len(r.recentOrders) {
		limit = len(r.recentOrders)
	}

	result := make([]OrderSummary, limit)
	copy(result, r.recentOrders[len(r.recentOrders)-limit:])
	return result
}

