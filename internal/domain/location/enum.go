package location

// ═══════════════════════════════════════════════════════════════
//              Schedule Type
// ═══════════════════════════════════════════════════════════════

type ScheduleType string

const (
	ScheduleTypeFixed ScheduleType = "fixed"
	ScheduleTypeMixed ScheduleType = "mixed"
)

func (s ScheduleType) IsValid() bool {
	switch s {
	case ScheduleTypeFixed, ScheduleTypeMixed:
		return true
	}
	return false
}

// String — для логов и дебага
func (s ScheduleType) String() string {
	return string(s)
}

// ParseScheduleType — парсинг из строки (API, БД)
func ParseScheduleType(s string) (ScheduleType, error) {
	st := ScheduleType(s)
	if !st.IsValid() {
		return "", ErrInvalidScheduleType
	}
	return st, nil
}
