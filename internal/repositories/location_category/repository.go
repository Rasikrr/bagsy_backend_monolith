package locationcategory

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type model struct {
	ID        uuid.UUID  `db:"id"`
	Slug      string     `db:"slug"`
	Name      string     `db:"name"`
	SortOrder int        `db:"sort_order"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

func (m model) toDomain() (*location.Category, error) {
	slug, err := shared.ParseSlug(m.Slug)
	if err != nil {
		return nil, fmt.Errorf("parse category slug: %w", err)
	}
	return &location.Category{
		ID:        m.ID,
		Slug:      slug,
		Name:      m.Name,
		SortOrder: m.SortOrder,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}

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

func (r *Repository) GetAll(ctx context.Context) ([]*location.Category, error) {
	var ms []model
	if err := pgxscan.Select(ctx, r.db, &ms, getAll); err != nil {
		return nil, fmt.Errorf("get all location categories: %w", err)
	}
	categories := make([]*location.Category, 0, len(ms))
	for _, m := range ms {
		cat, err := m.toDomain()
		if err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return categories, nil
}
