package media

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateMedia создает новую запись медиа в БД
func (r *Repository) CreateMedia(ctx context.Context, media *entity.Media) error {
	m := convert(media)

	_, err := r.db.Exec(ctx, createMediaSQL,
		m.ID,
		m.FileKey,
		m.BucketName,
		m.OriginalFilename,
		m.MimeType,
		m.SizeBytes,
		m.Width,
		m.Height,
		m.Status,
		m.UploadedBy,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to create media in db", err)
	}

	return nil
}

// GetMediaByID получает медиа по ID
func (r *Repository) GetMediaByID(ctx context.Context, id uuid.UUID) (*entity.Media, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getMediaByIDSQL, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("media not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get media from db", err)
	}
	out, convErr := m.convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get media from db", convErr)
	}
	return out, nil
}

// GetMediaByFileKey получает медиа по file_key
func (r *Repository) GetMediaByFileKey(ctx context.Context, fileKey string) (*entity.Media, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getMediaByFileKeySQL, fileKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("media not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get media from db", err)
	}
	out, convErr := m.convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get media from db", convErr)
	}
	return out, nil
}

// UpdateMediaStatus обновляет статус медиа
func (r *Repository) UpdateMediaStatus(ctx context.Context, id uuid.UUID, status enum.MediaStatus) error {
	result, err := r.db.Exec(ctx, updateMediaStatusSQL, id, status.String())
	if err != nil {
		return domainErr.NewInternalError("failed to update media status", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

// UpdateMediaMetadata обновляет метаданные медиа (width, height, size)
func (r *Repository) UpdateMediaMetadata(ctx context.Context, id uuid.UUID, width, height *int, sizeBytes int64) error {
	result, err := r.db.Exec(ctx, updateMediaMetadataSQL, id, width, height, sizeBytes)
	if err != nil {
		return domainErr.NewInternalError("failed to update media metadata", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

// SoftDeleteMedia помечает медиа как удаленное
func (r *Repository) SoftDeleteMedia(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, softDeleteMediaSQL, id)
	if err != nil {
		return domainErr.NewInternalError("failed to delete media", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

// SoftDeleteMediaByIDs помечает несколько медиа как удаленные
func (r *Repository) SoftDeleteMediaByIDs(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := r.db.Exec(ctx, softDeleteMediaByIDsSQL, ids)
	if err != nil {
		return domainErr.NewInternalError("failed to delete media files", err)
	}

	return nil
}

// ExistsMediaByFileKey проверяет существование медиа по file_key
func (r *Repository) ExistsMediaByFileKey(ctx context.Context, fileKey string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, existsByFileKeySQL, fileKey).Scan(&exists)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check media existence", err)
	}

	return exists, nil
}

// ListMediaByStatus получает список медиа по статусу (для cleanup jobs)
func (r *Repository) ListMediaByStatus(ctx context.Context, status enum.MediaStatus, limit int) ([]*entity.Media, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, listByStatusSQL, status.String(), limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Media{}, nil
		}
		return nil, domainErr.NewInternalError("failed to list media by status", err)
	}
	out, convErr := mm.convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to list media by status", convErr)
	}
	return out, nil
}
