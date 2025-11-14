package orders

import "github.com/bdv/gprs/internal/models/shared"

type Status string

type Order struct {
	ID            string
	UserID        string
	ProductID     string
	Price         shared.Money
	Status        Status
	PaymentID     string
	PaymentStatus string
}

const (
	StatusPending   Status = "PENDING"
	StatusCompleted Status = "COMPLETED"
	StatusFailed    Status = "FAILED"
)

func New(id, userID, productID string, price shared.Money) *Order {
	return &Order{
		ID:        id,
		UserID:    userID,
		ProductID: productID,
		Price:     price,
		Status:    StatusPending,
	}
}

func (o *Order) MarkPayment(paymentID, status string) {
	o.PaymentID = paymentID
	o.PaymentStatus = status
	if status == "AUTHORIZED" {
		o.Status = StatusCompleted
	} else {
		o.Status = StatusFailed
	}
}
