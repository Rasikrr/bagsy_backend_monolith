package media

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// TODO: review

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Asset
// ─────────────────────────────────────────────────────────────────

// Asset представляет метаданные физического файла в хранилище (S3).
type Asset struct {
	ID        uuid.UUID
	Bucket    string
	ObjectKey string
	Filename  string
	MimeType  string
	SizeBytes int64
	Status    Status
	CreatedAt time.Time
	UpdatedAt *time.Time
}

// ─────────────────────────────────────────────────────────────────
// Factory
// ─────────────────────────────────────────────────────────────────

type CreateAssetParams struct {
	Bucket    string
	ObjectKey string
	Filename  string
	MimeType  string
	SizeBytes int64
}

// NewAsset создает новую запись о медиафайле.
// По умолчанию статус всегда Pending, так как физическая загрузка файла
// происходит асинхронно с фронтенда напрямую в S3.
func NewAsset(params CreateAssetParams) (*Asset, error) {
	if params.SizeBytes <= 0 {
		return nil, ErrInvalidFileSize
	}

	cleanFilename := strings.TrimSpace(params.Filename)
	if cleanFilename == "" {
		return nil, ErrEmptyFilename
	}

	cleanMime := strings.TrimSpace(params.MimeType)
	if cleanMime == "" {
		return nil, ErrEmptyMimeType
	}

	return &Asset{
		ID:        uuid.New(),
		Bucket:    params.Bucket,
		ObjectKey: params.ObjectKey,
		Filename:  cleanFilename,
		MimeType:  cleanMime,
		SizeBytes: params.SizeBytes,
		Status:    StatusPending, // Начальное состояние
		CreatedAt: time.Now(),
	}, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods (State Transitions)
// ─────────────────────────────────────────────────────────────────

// MarkAsUploaded подтверждает, что фронтенд успешно загрузил байты в S3.
func (a *Asset) MarkAsUploaded() {
	if a.Status == StatusUploaded {
		return
	}
	a.Status = StatusUploaded
	a.touch()
}

// MarkAsFailed помечает загрузку как неуспешную (например, истек таймаут Presigned URL).
func (a *Asset) MarkAsFailed() {
	if a.Status == StatusFailed {
		return
	}
	a.Status = StatusFailed
	a.touch()
}

// IsReady проверяет, можно ли использовать этот файл (например, прикреплять к локации).
func (a *Asset) IsReady() bool {
	return a.Status == StatusUploaded
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (a *Asset) touch() {
	now := time.Now()
	a.UpdatedAt = &now
}
