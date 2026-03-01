package media

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*media.Asset, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, media.ErrAssetNotFound
		}
		return nil, fmt.Errorf("get media asset by id: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*media.Asset, error) {
	var models []model
	if err := pgxscan.Select(ctx, r.db, &models, getByIDs, ids); err != nil {
		return nil, fmt.Errorf("get media assets by ids: %w", err)
	}

	result := make([]*media.Asset, 0, len(models))
	for _, m := range models {
		asset, err := m.toDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, asset)
	}
	return result, nil
}

func (r *Repository) Save(ctx context.Context, asset *media.Asset) error {
	m := fromDomain(asset)
	_, err := r.db.Exec(ctx, saveAsset,
		m.ID,
		m.Bucket,
		m.ObjectKey,
		m.Filename,
		m.MimeType,
		m.SizeBytes,
		m.Status,
		m.CreatedAt,
		m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save media asset: %w", err)
	}
	return nil
}

func (r *Repository) MarkExpiredPendingAsFailed(ctx context.Context, threshold time.Time) (int64, error) {
	res, err := r.db.Exec(ctx, markExpiredPendingAsFailed, threshold)
	if err != nil {
		return 0, fmt.Errorf("mark expired pending media as failed: %w", err)
	}
	return res.RowsAffected(), nil
}
