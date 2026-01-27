package notification

//go:generate enumer -type=Type -json -trimprefix Type -transform=snake -output type_enumer.go
type Type int8

const (
	TypeDayBefore  Type = iota // За день до записи
	TypeHourBefore             // За час до записи
)

// GetOffset возвращает смещение времени для данного типа уведомления
func (t Type) GetOffset() int {
	switch t {
	case TypeDayBefore:
		return 24 * 60 // минуты
	case TypeHourBefore:
		return 60 // минуты
	default:
		return 0
	}
}

// AllTypes возвращает все типы уведомлений
func AllTypes() []Type {
	return []Type{TypeDayBefore, TypeHourBefore}
}
