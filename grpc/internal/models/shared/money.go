package shared

type Money struct {
	Currency string
	Amount   float64
}

func NewMoney(currency string, amount float64) Money {
	return Money{Currency: currency, Amount: amount}
}

