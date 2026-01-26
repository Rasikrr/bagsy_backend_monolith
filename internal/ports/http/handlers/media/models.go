package media

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"time"
)

//go:generate easyjson -all models.go

type getUploadURLRequest struct {
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content_type" validate:"required"`
	Purpose     string `json:"purpose" validate:"required"`
}

func (r *getUploadURLRequest) Validate() error {
	if err := request.GetValidator().Struct(r); err != nil {
		return request.HandleValidationError(err)
	}
	return nil
}

type getUploadURLResponse struct {
	MediaID   string `json:"media_id"`
	URL       string `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

func toUploadURLResponse(dto *media.UploadedMedia) *getUploadURLResponse {
	return &getUploadURLResponse{
		MediaID:   dto.MediaID.String(),
		URL:       dto.URL,
		ExpiresAt: dto.ExpiresAt,
	}
}
