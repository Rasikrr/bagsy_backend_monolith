package service_categories

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type Repository interface {
	GetByID(ctx context.Context, id int) (*entity.ServiceCategory, error)
	Create(ctx context.Context, category *entity.ServiceCategory) error
	Update(ctx context.Context, category *entity.ServiceCategory) error
	Delete(ctx context.Context, categories ...*entity.ServiceCategory) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db}
}

func (r *repository) GetByID(ctx context.Context, id int) (*entity.ServiceCategory, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getServiceCategoryByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrServiceCategoryNotFound
		}
		return nil, err
	}
	return m.convert(), nil
}

func (r *repository) Create(ctx context.Context, category *entity.ServiceCategory) error {
	m := convert(category)
	err := r.db.QueryRow(ctx, createServiceCategory, m.Name, m.Description, m.UpdatedBy).Scan(&category.ID)
	return err
}

func (r *repository) Update(ctx context.Context, category *entity.ServiceCategory) error {
	m := convert(category)
	_, err := r.db.Exec(ctx, updateServiceCategory, m.ID, m.Name, m.Description, m.UpdatedBy)
	return err
}

func (r *repository) Delete(ctx context.Context, categories ...*entity.ServiceCategory) error {
	ids := lo.Map(categories, func(item *entity.ServiceCategory, _ int) int {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteServiceCategory, pq.Array(ids))
	return err
}