package portfolio

type Money struct {
	Amount   float64
	Currency string
}

func NewMoney(amount float64, currency string) Money {
	return Money{Amount: amount, Currency: currency}
}
