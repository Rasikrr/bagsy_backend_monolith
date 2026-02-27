package media

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/google/uuid"
)

type model struct {
	ID        uuid.UUID  `db:"id"`
	Bucket    string     `db:"bucket"`
	ObjectKey string     `db:"object_key"`
	Filename  string     `db:"filename"`
	MimeType  string     `db:"mime_type"`
	SizeBytes int64      `db:"size_bytes"`
	Status    string     `db:"status"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func fromDomain(a *media.Asset) *model {
	return &model{
		ID:        a.ID,
		Bucket:    a.Bucket,
		ObjectKey: a.ObjectKey,
		Filename:  a.Filename,
		MimeType:  a.MimeType.String(),
		SizeBytes: a.SizeBytes,
		Status:    string(a.Status),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (m *model) toDomain() (*media.Asset, error) {
	mimeType, err := media.ParseMimeType(m.MimeType)
	if err != nil {
		return nil, err
	}

	return &media.Asset{
		ID:        m.ID,
		Bucket:    m.Bucket,
		ObjectKey: m.ObjectKey,
		Filename:  m.Filename,
		MimeType:  mimeType,
		SizeBytes: m.SizeBytes,
		Status:    media.Status(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
