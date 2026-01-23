package notification

import (
	"time"

	"github.com/google/uuid"
)

// Notification представляет уведомление о предстоящей записи
type Notification struct {
	ID            uuid.UUID
	BagsyID       uuid.UUID
	Type          Type
	RecipientType RecipientType
	ScheduledAt   time.Time  // Когда должно быть отправлено
	SentAt        *time.Time // Когда фактически отправлено
	Status        Status
	Attempts      int
	LastError     *string
	CreatedAt     time.Time
}
