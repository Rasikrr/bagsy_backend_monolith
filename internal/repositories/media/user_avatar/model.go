package useravatar

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/google/uuid"
)

// model представляет DB модель связи пользователя с аватаром
type model struct {
	UserPhone string     `db:"user_phone"`
	MediaID   uuid.UUID  `db:"media_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// convert преобразует entity.UserMedia → DB model
func convert(e *entity.UserMedia) model {
	return model{
		UserPhone: e.UserPhone,
		MediaID:   e.MediaID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// convert преобразует DB model → entity.UserMedia
func (m model) convert() *entity.UserMedia {
	return &entity.UserMedia{
		UserPhone: m.UserPhone,
		MediaID:   m.MediaID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
