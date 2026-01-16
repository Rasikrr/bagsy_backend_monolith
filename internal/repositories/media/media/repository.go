package media

import (
	"context"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/models"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

// Repository отвечает за базовую работу с media таблицей
type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// Create создает новую запись медиа в БД
func (r *Repository) Create(ctx context.Context, media *media.Media) error {
	m := models.FromEntity(media)

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

// GetByID получает медиа по ID
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*media.Media, error) {
	var m models.Media
	err := pgxscan.Get(ctx, r.db, &m, getMediaByIDSQL, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("media not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get media from db", err)
	}
	out, convErr := m.Convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get media from db", convErr)
	}
	return out, nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids ...uuid.UUID) ([]*media.Media, error) {
	var mm models.MediaList
	err := pgxscan.Select(ctx, r.db, &mm, getMediaByIDsSQL, pq.Array(ids))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}
	return mm.Convert()
}

// GetByFileKey получает медиа по file_key
func (r *Repository) GetByFileKey(ctx context.Context, fileKey string) (*media.Media, error) {
	var m models.Media
	err := pgxscan.Get(ctx, r.db, &m, getMediaByFileKeySQL, fileKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("media not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get media from db", err)
	}
	out, convErr := m.Convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get media from db", convErr)
	}
	return out, nil
}

// UpdateStatus обновляет статус медиа
func (r *Repository) UpdateStatus(ctx context.Context, id uuid.UUID, status media.Status) error {
	result, err := r.db.Exec(ctx, updateMediaStatusSQL, id, status.String())
	if err != nil {
		return domainErr.NewInternalError("failed to update media status", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

func (r *Repository) UpdateStatuses(ctx context.Context, ids []uuid.UUID, status media.Status) error {
	result, err := r.db.Exec(ctx, updateMediaStatusesSQL, pq.Array(ids), status.String())
	if err != nil {
		return domainErr.NewInternalError("failed to update media status", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

// UpdateMetadata обновляет метаданные медиа (width, height, size)
func (r *Repository) UpdateMetadata(ctx context.Context, id uuid.UUID, width, height *int, sizeBytes int64) error {
	result, err := r.db.Exec(ctx, updateMediaMetadataSQL, id, width, height, sizeBytes)
	if err != nil {
		return domainErr.NewInternalError("failed to update media metadata", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

// SoftDeleteByID помечает медиа как удаленное
func (r *Repository) SoftDeleteByID(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx, softDeleteMediaSQL, id)
	if err != nil {
		return domainErr.NewInternalError("failed to delete media", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("media not found", nil)
	}

	return nil
}

// SoftDeleteByIDs помечает несколько медиа как удаленные
func (r *Repository) SoftDeleteByIDs(ctx context.Context, ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := r.db.Exec(ctx, softDeleteMediaByIDsSQL, ids)
	if err != nil {
		return domainErr.NewInternalError("failed to delete media files", err)
	}

	return nil
}

// ExistsByFileKey проверяет существование медиа по file_key
func (r *Repository) ExistsByFileKey(ctx context.Context, fileKey string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, existsByFileKeySQL, fileKey).Scan(&exists)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check media existence", err)
	}

	return exists, nil
}

// ListByStatus получает список медиа по статусу (для cleanup jobs)
func (r *Repository) ListByStatus(ctx context.Context, status media.Status, limit int) ([]*media.Media, error) {
	var mm models.MediaList
	err := pgxscan.Select(ctx, r.db, &mm, listByStatusSQL, status.String(), limit)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*media.Media{}, nil
		}
		return nil, domainErr.NewInternalError("failed to list media by status", err)
	}
	out, convErr := mm.Convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to list media by status", convErr)
	}
	return out, nil
}
