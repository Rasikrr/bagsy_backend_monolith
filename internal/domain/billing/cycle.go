package billing

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
