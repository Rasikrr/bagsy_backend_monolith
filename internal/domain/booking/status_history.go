package booking

import (
	"time"

	"github.com/google/uuid"
)

// StatusHistoryEntry — запись в истории статусов (Entity внутри Aggregate)
type StatusHistoryEntry struct {
	ID         uuid.UUID
	FromStatus *Status
	ToStatus   Status
	ChangedBy  *uuid.UUID
	Reason     *string
	CreatedAt  time.Time
}
