package pointmedia

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/google/uuid"
)

// model представляет DB модель связи точки с её фотографиями
type model struct {
	ID           uuid.UUID  `db:"id"`
	PointCode    string     `db:"point_code"`
	MediaID      uuid.UUID  `db:"media_id"`
	DisplayOrder int        `db:"display_order"`
	CreatedAt    time.Time  `db:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"`
	DeletedAt    *time.Time `db:"deleted_at"`
}

type modelList []model

func (mm modelList) convert() []*entity.PointMedia {
	out := make([]*entity.PointMedia, len(mm))
	for i, m := range mm {
		out[i] = m.convert()
	}
	return out
}

// convert преобразует entity.PointMedia → DB model
func convert(e *entity.PointMedia) model {
	m := model{
		ID:           e.ID,
		PointCode:    e.PointCode,
		MediaID:      e.MediaID,
		DisplayOrder: e.DisplayOrder,
		CreatedAt:    e.CreatedAt,
		DeletedAt:    e.DeletedAt,
	}

	// UpdatedAt - если не nil, используем значение, иначе NOW()
	if e.UpdatedAt != nil {
		m.UpdatedAt = *e.UpdatedAt
	} else {
		m.UpdatedAt = e.CreatedAt
	}

	return m
}

// convert преобразует DB model → entity.PointMedia
func (m model) convert() *entity.PointMedia {
	return &entity.PointMedia{
		ID:           m.ID,
		PointCode:    m.PointCode,
		MediaID:      m.MediaID,
		DisplayOrder: m.DisplayOrder,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    &m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
	}
}
