package notification

//go:generate enumer -type=Status -json -trimprefix Status -transform=snake -output status_enumer.go
type Status int8

const (
	StatusPending Status = iota // Ожидает отправки
	StatusSent                  // Отправлено
	StatusFailed                // Ошибка отправки
	StatusSkipped               // Пропущено (например, время уже прошло)
)
