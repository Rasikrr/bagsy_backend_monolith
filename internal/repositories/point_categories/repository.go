package pointcategories

import (
	"context"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int) (*point.Category, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getPointCategoryByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, point.ErrPointCategoryNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point category from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) ExistsByID(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := pgxscan.Get(ctx, r.db, &exists, existsByID, id)
	if err != nil {
		return false, domainErr.NewInternalError("failed to get point category from db", err)
	}
	return exists, nil
}

func (r *Repository) Create(ctx context.Context, category *point.Category) error {
	m := convert(category)
	err := r.db.QueryRow(ctx, createPointCategory, m.Name, m.Description, m.UpdatedBy).Scan(&category.ID)
	if err != nil {
		return domainErr.NewInternalError("failed to create point category in db", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, category *point.Category) error {
	m := convert(category)
	_, err := r.db.Exec(ctx, updatePointCategory, m.ID, m.Name, m.Description, m.UpdatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to update point category in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, categories ...*point.Category) error {
	ids := lo.Map(categories, func(item *point.Category, _ int) int {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deletePointCategory, pq.Array(ids))
	if err != nil {
		return domainErr.NewInternalError("failed to delete point categories from db", err)
	}
	return nil
}

func (r *Repository) GetAll(ctx context.Context) ([]*point.Category, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getAllPointCategories)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*point.Category{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get point categories from db", err)
	}
	return mm.convert(), nil
}
