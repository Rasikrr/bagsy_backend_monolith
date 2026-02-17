package billing

import "time"

type Cycle string

const (
	CycleMonthly Cycle = "monthly"
	CycleAnnual  Cycle = "annual"
)

func (c Cycle) IsValid() bool {
	switch c {
	case CycleMonthly, CycleAnnual:
		return true
	}
	return false
}

func (c Cycle) Duration() time.Duration {
	switch c {
	case CycleAnnual:
		return 365 * 24 * time.Hour
	default:
		return 30 * 24 * time.Hour
	}
}
