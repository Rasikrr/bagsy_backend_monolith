package media

import (
	"time"

	"github.com/google/uuid"
)

type UploadedMedia struct {
	MediaID   uuid.UUID
	URL       string
	ExpiresAt time.Time
}
