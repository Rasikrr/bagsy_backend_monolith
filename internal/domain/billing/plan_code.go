package billing

type PlanCode string

const (
	PlanCodeSolo    PlanCode = "solo"
	PlanCodePoint   PlanCode = "point"
	PlanCodeNetwork PlanCode = "network"
)

func (c PlanCode) IsValid() bool {
	switch c {
	case PlanCodeSolo, PlanCodePoint, PlanCodeNetwork:
		return true
	}
	return false
}

func (c PlanCode) IsSolo() bool {
	return c == PlanCodeSolo
}

func (c PlanCode) IsPoint() bool {
	return c == PlanCodePoint
}

func (c PlanCode) IsNetwork() bool {
	return c == PlanCodeNetwork
}

func (c PlanCode) String() string {
	return string(c)
}

// TrialDays returns the trial duration for the plan.
// Solo = 60 days (2 months), Point/Network = 30 days (1 month).
func (c PlanCode) TrialDays() int {
	if c == PlanCodeSolo {
		return 60
	}
	return 30
}

func ParsePlanCode(s string) (PlanCode, error) {
	plan := PlanCode(s)
	if !plan.IsValid() {
		return PlanCodeSolo, ErrInvalidPlanCode
	}
	return plan, nil
}
