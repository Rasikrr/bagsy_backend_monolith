// nolint
package pointcategoryservices

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

func (r *Repository) GetByID(ctx context.Context, id int) (*point.CategoryService, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, point.ErrPointCategoryServiceNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point category service from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) GetByPointCategoryID(ctx context.Context, pointCategoryID int) ([]*point.CategoryService, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByPointCategoryID, pointCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, point.ErrPointCategoryServiceNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point category services from db", err)
	}
	return mm.convert(), nil
}

func (r *Repository) GetByServiceCategoryID(ctx context.Context, serviceCategoryID int) ([]*point.CategoryService, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByServiceCategoryID, serviceCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, point.ErrPointCategoryServiceNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point category services from db", err)
	}
	return mm.convert(), nil
}

func (r *Repository) GetByPointCategoryIDAndServiceCategoryID(ctx context.Context, pointCategoryID, serviceCategoryID int) (*point.CategoryService, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getByPointCategoryIDAndServiceCategoryID, pointCategoryID, serviceCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, point.ErrPointCategoryServiceNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point category service from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, pcs *point.CategoryService) error {
	m := convert(pcs)
	err := r.db.QueryRow(ctx, createPointCategoryService, m.PointCategoryID, m.ServiceCategoryID).Scan(&pcs.ID)
	if err != nil {
		return domainErr.NewInternalError("failed to create point category service in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, pointCategoryServices ...*point.CategoryService) error {
	ids := lo.Map(pointCategoryServices, func(item *point.CategoryService, _ int) int {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deletePointCategoryService, pq.Array(ids))
	if err != nil {
		return domainErr.NewInternalError("failed to delete point category services from db", err)
	}
	return nil
}

func (r *Repository) DeleteByPointCategoryID(ctx context.Context, pointCategoryID int) error {
	_, err := r.db.Exec(ctx, deleteByPointCategoryID, pointCategoryID)
	if err != nil {
		return domainErr.NewInternalError("failed to delete point category services from db", err)
	}
	return nil
}

func (r *Repository) DeleteByPointCategoryIDAndServiceCategoryIDs(ctx context.Context, pointCategoryID int, serviceCategoryIDs []int) error {
	if len(serviceCategoryIDs) == 0 {
		return nil
	}
	_, err := r.db.Exec(ctx, deleteByPointCategoryIDAndServiceCategoryIDs, pointCategoryID, pq.Array(serviceCategoryIDs))
	if err != nil {
		return domainErr.NewInternalError("failed to delete point category services from db", err)
	}
	return nil
}
