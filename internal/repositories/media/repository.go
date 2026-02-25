package media

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/core/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
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
