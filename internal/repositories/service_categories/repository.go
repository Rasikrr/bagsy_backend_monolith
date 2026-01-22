package servicecategories

import (
	"context"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
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

func (r *Repository) GetByID(ctx context.Context, id int) (*service.Category, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getServiceCategoryByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrServiceCategoryNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get service category from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, category *service.Category) error {
	m := convert(category)
	err := r.db.QueryRow(ctx, createServiceCategory, m.Name, m.Description, m.UpdatedBy).Scan(&category.ID)
	if err != nil {
		return domainErr.NewInternalError("failed to create service category in db", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, category *service.Category) error {
	m := convert(category)
	_, err := r.db.Exec(ctx, updateServiceCategory, m.ID, m.Name, m.Description, m.UpdatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to update service category in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, categories ...*service.Category) error {
	ids := lo.Map(categories, func(item *service.Category, _ int) int {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteServiceCategory, pq.Array(ids))
	if err != nil {
		return domainErr.NewInternalError("failed to delete service categories from db", err)
	}
	return nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []int) ([]*service.Category, error) {
	if len(ids) == 0 {
		return []*service.Category{}, nil
	}
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getServiceCategoriesByIDs, pq.Array(ids))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*service.Category{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get service categories from db", err)
	}
	return mm.convert(), nil
}
