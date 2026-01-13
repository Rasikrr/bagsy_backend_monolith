package pointmedia

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/media/models"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Repository отвечает за работу с point_media таблицей
type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

// Add добавляет фото к точке
func (r *Repository) Add(ctx context.Context, pointMedia *entity.PointMedia) error {
	m := convert(pointMedia)

	_, err := r.db.Exec(ctx, addPointPhotoSQL,
		m.ID,
		m.PointCode,
		m.MediaID,
		m.DisplayOrder,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to add point photo", err)
	}

	return nil
}

// GetAll получает все фото точки (только связи, без Media)
func (r *Repository) GetAll(ctx context.Context, pointCode string) ([]*entity.PointMedia, error) {
	var mm modelList
	err := pgxscan.Select(ctx, r.db, &mm, getPointPhotosSQL, pointCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.PointMedia{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point photos", err)
	}

	return mm.convert(), nil
}

// GetWithMedia получает все фото точки с полными данными Media через JOIN
// Использует эффективный SQL JOIN вместо N+1 запросов
func (r *Repository) GetWithMedia(ctx context.Context, pointCode string) ([]*entity.Media, error) {
	var mm models.MediaList
	err := pgxscan.Select(ctx, r.db, &mm, getPointPhotosWithMediaSQL, pointCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Media{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point photos with media", err)
	}

	out, convErr := mm.Convert()
	if convErr != nil {
		return nil, domainErr.NewInternalError("failed to get point photos with media", convErr)
	}

	return out, nil
}

// Get получает одну связь точки с фото
func (r *Repository) Get(ctx context.Context, pointCode string, mediaID uuid.UUID) (*entity.PointMedia, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getPointPhotoSQL, pointCode, mediaID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("point photo not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get point photo", err)
	}

	return m.convert(), nil
}

// UpdateOrder обновляет порядок отображения фото
func (r *Repository) UpdateOrder(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error {
	result, err := r.db.Exec(ctx, updatePointPhotoOrderSQL, pointCode, mediaID, displayOrder)
	if err != nil {
		return domainErr.NewInternalError("failed to update point photo order", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("point photo not found", nil)
	}

	return nil
}

// Remove удаляет фото точки (soft delete)
func (r *Repository) Remove(ctx context.Context, pointCode string, mediaID uuid.UUID) error {
	result, err := r.db.Exec(ctx, removePointPhotoSQL, pointCode, mediaID)
	if err != nil {
		return domainErr.NewInternalError("failed to remove point photo", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("point photo not found", nil)
	}

	return nil
}

// RemoveAll удаляет все фото точки (soft delete)
func (r *Repository) RemoveAll(ctx context.Context, pointCode string) error {
	_, err := r.db.Exec(ctx, removeAllPointPhotosSQL, pointCode)
	if err != nil {
		return domainErr.NewInternalError("failed to remove all point photos", err)
	}

	return nil
}

// Count подсчитывает количество фото у точки
func (r *Repository) Count(ctx context.Context, pointCode string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, countPointPhotosSQL, pointCode).Scan(&count)
	if err != nil {
		return 0, domainErr.NewInternalError("failed to count point photos", err)
	}

	return count, nil
}

// Has проверяет, есть ли у точки указанное фото
func (r *Repository) Has(ctx context.Context, pointCode string, mediaID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, pointHasPhotoSQL, pointCode, mediaID).Scan(&exists)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check point photo existence", err)
	}

	return exists, nil
}
