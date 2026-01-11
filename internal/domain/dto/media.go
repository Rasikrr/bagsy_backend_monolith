package dto

import (
	"time"

	"github.com/google/uuid"
)

type UploadMediaResponse struct {
	MediaID   uuid.UUID
	URL       string
	ExpiresAt time.Time
}
