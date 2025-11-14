package payments

type Status string

type Payment struct {
	ID       string
	OrderID  string
	Amount   float64
	Currency string
	Status   Status
	Message  string
}

const (
	StatusAuthorized Status = "AUTHORIZED"
	StatusDeclined   Status = "DECLINED"
)

func New(id, orderID string, amount float64, currency string, status Status, message string) *Payment {
	return &Payment{
		ID:       id,
		OrderID:  orderID,
		Amount:   amount,
		Currency: currency,
		Status:   status,
		Message:  message,
	}
}

func (p *Payment) UpdateStatus(status Status, message string) {
	p.Status = status
	p.Message = message
}
