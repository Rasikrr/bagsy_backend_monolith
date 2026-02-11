package shared

import "github.com/shopspring/decimal"

type Money struct {
	amount decimal.Decimal
}

func NewMoney(amount decimal.Decimal) (Money, error) {
	if amount.IsNegative() {
		return Money{}, ErrNegativeAmount
	}
	return Money{amount: amount}, nil
}

func NewMoneyFromInt(amount int64) (Money, error) {
	return NewMoney(decimal.NewFromInt(amount))
}

func NewMoneyFromFloat(amount float64) (Money, error) {
	return NewMoney(decimal.NewFromFloat(amount))
}

func (m Money) IsZero() bool {
	return m.amount.IsZero()
}

func (m Money) String() string {
	return m.amount.StringFixed(2)
}

func (m Money) Amount() decimal.Decimal {
	return m.amount
}
