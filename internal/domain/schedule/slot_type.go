package schedule

type SlotType string

const (
	SlotTypeWork SlotType = "work"
	SlotTypeRest SlotType = "rest"
)

func (s SlotType) IsValid() bool {
	switch s {
	case SlotTypeWork, SlotTypeRest:
		return true
	}
	return false
}

func (s SlotType) String() string {
	return string(s)
}

func ParseSlotType(s string) (SlotType, error) {
	st := SlotType(s)
	if !st.IsValid() {
		return "", ErrInvalidSlotType
	}
	return st, nil
}
