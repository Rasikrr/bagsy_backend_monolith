package models

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// Media представляет DB модель медиа-файла с тегами
type Media struct {
	ID               uuid.UUID  `db:"id"`
	FileKey          string     `db:"file_key"`
	BucketName       string     `db:"bucket_name"`
	OriginalFilename string     `db:"original_filename"`
	MimeType         string     `db:"mime_type"`
	SizeBytes        int64      `db:"size_bytes"`
	Width            *int       `db:"width"`
	Height           *int       `db:"height"`
	Status           string     `db:"status"`
	UploadedBy       *string    `db:"uploaded_by"`
	CreatedAt        time.Time  `db:"created_at"`
	UpdatedAt        time.Time  `db:"updated_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
}

type MediaList []Media

func (mm MediaList) Convert() ([]*media.Media, error) {
	out := make([]*media.Media, len(mm))
	for i, m := range mm {
		e, err := m.Convert()
		if err != nil {
			return nil, err
		}
		out[i] = e
	}
	return out, nil
}

// FromEntity преобразует entity.Media → DB model
func FromEntity(e *media.Media) Media {
	m := Media{
		ID:               e.ID,
		FileKey:          e.FileKey,
		BucketName:       e.BucketName,
		OriginalFilename: e.OriginalFilename,
		MimeType:         e.MimeType,
		SizeBytes:        e.SizeBytes,
		Width:            e.Width,
		Height:           e.Height,
		Status:           e.Status.String(),
		UploadedBy:       e.UploadedBy,
		CreatedAt:        e.CreatedAt,
		DeletedAt:        e.DeletedAt,
	}

	// UpdatedAt - если не nil, используем значение, иначе NOW()
	if e.UpdatedAt != nil {
		m.UpdatedAt = *e.UpdatedAt
	} else {
		m.UpdatedAt = e.CreatedAt
	}

	return m
}

// Convert преобразует DB model → entity.Media
func (m Media) Convert() (*media.Media, error) {
	status, err := media.StatusString(m.Status)
	if err != nil {
		return nil, errors.Wrap(err, "invalid media status")
	}

	return &media.Media{
		ID:               m.ID,
		FileKey:          m.FileKey,
		BucketName:       m.BucketName,
		OriginalFilename: m.OriginalFilename,
		MimeType:         m.MimeType,
		SizeBytes:        m.SizeBytes,
		Width:            m.Width,
		Height:           m.Height,
		Status:           status,
		UploadedBy:       m.UploadedBy,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        &m.UpdatedAt,
		DeletedAt:        m.DeletedAt,
	}, nil
}
