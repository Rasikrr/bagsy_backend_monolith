package billing

type Limit struct {
	value     int
	unlimited bool
}

func NewLimit(val int) (Limit, error) {
	if val < 0 {
		return Limit{}, ErrNegativeLimit
	}
	return Limit{
		value:     val,
		unlimited: false,
	}, nil
}

func NewUnlimited() Limit {
	return Limit{
		unlimited: true,
	}
}

// IsExceeded returns true if current count reaches or exceeds the limit.
// Unlimited resources can never be exceeded.
func (l Limit) IsExceeded(current int) bool {
	if l.unlimited {
		return false
	}
	return current >= l.value
}

func (l Limit) IsUnlimited() bool {
	return l.unlimited
}

func (l Limit) Value() int {
	return l.value
}
