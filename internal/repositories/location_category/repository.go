package locationcategory

import (
	"context"
	"fmt"

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

func (r *Repository) ExistsByID(ctx context.Context, id uuid.UUID) (bool, error) {
	var exists bool
	if err := pgxscan.Get(ctx, r.db, &exists, existsByID, id); err != nil {
		return false, fmt.Errorf("check location_category exists: %w", err)
	}
	return exists, nil
}
