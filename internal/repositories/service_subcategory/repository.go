package servicesubcategory

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
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

func (r *Repository) GetByID(ctx context.Context, id int) (*entity.ServiceSubcategory, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getServiceSubcategoryByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrServiceSubcategoryNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get service subcategory from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, subcategory *entity.ServiceSubcategory) error {
	m := convert(subcategory)
	err := r.db.QueryRow(ctx, createServiceSubcategory, m.ServiceCategoryID, m.Name, m.Description, m.UpdatedBy).Scan(&subcategory.ID)
	if err != nil {
		return domainErr.NewInternalError("failed to create service subcategory in db", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, subcategory *entity.ServiceSubcategory) error {
	m := convert(subcategory)
	_, err := r.db.Exec(ctx, updateServiceSubcategory, m.ID, m.ServiceCategoryID, m.Name, m.Description, m.UpdatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to update service subcategory in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, subcategories ...*entity.ServiceSubcategory) error {
	ids := lo.Map(subcategories, func(item *entity.ServiceSubcategory, _ int) int {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteServiceSubcategory, pq.Array(ids))
	if err != nil {
		return domainErr.NewInternalError("failed to delete service subcategories from db", err)
	}
	return nil
}
