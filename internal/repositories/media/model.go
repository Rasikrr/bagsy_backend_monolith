package media

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

// model представляет DB модель медиа-файла с тегами
type model struct {
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

type models []model

func (mm models) convert() ([]*entity.Media, error) {
	out := make([]*entity.Media, len(mm))
	for i, m := range mm {
		e, err := m.convert()
		if err != nil {
			return nil, err
		}
		out[i] = e
	}
	return out, nil
}

// convert преобразует entity.Media → DB model
func convert(e *entity.Media) model {
	m := model{
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

// convert преобразует DB model → entity.Media
func (m model) convert() (*entity.Media, error) {
	status, err := enum.MediaStatusString(m.Status)
	if err != nil {
		return nil, errors.Wrap(err, "invalid media status")
	}

	return &entity.Media{
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

// userMediaModel представляет DB модель связи пользователя с аватаром
type userMediaModel struct {
	UserPhone string     `db:"user_phone"`
	MediaID   uuid.UUID  `db:"media_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// convert преобразует entity.UserMedia → DB model
func convertUserMedia(e *entity.UserMedia) userMediaModel {
	return userMediaModel{
		UserPhone: e.UserPhone,
		MediaID:   e.MediaID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}
}

// convert преобразует DB model → entity.UserMedia
func (m userMediaModel) convert() *entity.UserMedia {
	return &entity.UserMedia{
		UserPhone: m.UserPhone,
		MediaID:   m.MediaID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
