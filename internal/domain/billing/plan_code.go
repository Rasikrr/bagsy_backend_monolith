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
