package media

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// AddPointPhoto добавляет фото к точке
func (r *Repository) AddPointPhoto(ctx context.Context, pointMedia *entity.PointMedia) error {
	m := convertPointMedia(pointMedia)

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

// GetPointPhotos получает все фото точки (только связи, без Media)
func (r *Repository) GetPointPhotos(ctx context.Context, pointCode string) ([]*entity.PointMedia, error) {
	var mm pointMediaModels
	err := pgxscan.Select(ctx, r.db, &mm, getPointPhotosSQL, pointCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.PointMedia{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point photos", err)
	}

	return mm.convert(), nil
}

// GetPointPhotosWithMedia получает все фото точки с полными данными Media через JOIN
func (r *Repository) GetPointPhotosWithMedia(ctx context.Context, pointCode string) ([]*entity.Media, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getPointPhotosWithMediaSQL, pointCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*entity.Media{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point photos with media", err)
	}

	return mm.convert()
}

// GetPointPhoto получает одну связь точки с фото
func (r *Repository) GetPointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID) (*entity.PointMedia, error) {
	var m pointMediaModel
	err := pgxscan.Get(ctx, r.db, &m, getPointPhotoSQL, pointCode, mediaID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.NewNotFoundError("point photo not found", err)
		}
		return nil, domainErr.NewInternalError("failed to get point photo", err)
	}

	return m.convert(), nil
}

// UpdatePointPhotoOrder обновляет порядок отображения фото
func (r *Repository) UpdatePointPhotoOrder(ctx context.Context, pointCode string, mediaID uuid.UUID, displayOrder int) error {
	result, err := r.db.Exec(ctx, updatePointPhotoOrderSQL, pointCode, mediaID, displayOrder)
	if err != nil {
		return domainErr.NewInternalError("failed to update point photo order", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("point photo not found", nil)
	}

	return nil
}

// RemovePointPhoto удаляет фото точки (soft delete)
func (r *Repository) RemovePointPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID) error {
	result, err := r.db.Exec(ctx, removePointPhotoSQL, pointCode, mediaID)
	if err != nil {
		return domainErr.NewInternalError("failed to remove point photo", err)
	}

	if result.RowsAffected() == 0 {
		return domainErr.NewNotFoundError("point photo not found", nil)
	}

	return nil
}

// RemoveAllPointPhotos удаляет все фото точки (soft delete)
func (r *Repository) RemoveAllPointPhotos(ctx context.Context, pointCode string) error {
	_, err := r.db.Exec(ctx, removeAllPointPhotosSQL, pointCode)
	if err != nil {
		return domainErr.NewInternalError("failed to remove all point photos", err)
	}

	return nil
}

// CountPointPhotos подсчитывает количество фото у точки
func (r *Repository) CountPointPhotos(ctx context.Context, pointCode string) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, countPointPhotosSQL, pointCode).Scan(&count)
	if err != nil {
		return 0, domainErr.NewInternalError("failed to count point photos", err)
	}

	return count, nil
}

// PointHasPhoto проверяет, есть ли у точки указанное фото
func (r *Repository) PointHasPhoto(ctx context.Context, pointCode string, mediaID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, pointHasPhotoSQL, pointCode, mediaID).Scan(&exists)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check point photo existence", err)
	}

	return exists, nil
}
