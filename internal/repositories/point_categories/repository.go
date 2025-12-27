package point_categories

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

type Repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id int) (*entity.PointCategory, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getPointCategoryByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrPointCategoryNotFound.WithError(err)
		}
		return nil, err
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, category *entity.PointCategory) error {
	m := convert(category)
	err := r.db.QueryRow(ctx, createPointCategory, m.Name, m.Description, m.UpdatedBy).Scan(&category.ID)
	return err
}

func (r *Repository) Update(ctx context.Context, category *entity.PointCategory) error {
	m := convert(category)
	_, err := r.db.Exec(ctx, updatePointCategory, m.ID, m.Name, m.Description, m.UpdatedBy)
	return err
}

func (r *Repository) Delete(ctx context.Context, categories ...*entity.PointCategory) error {
	ids := lo.Map(categories, func(item *entity.PointCategory, _ int) int {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deletePointCategory, pq.Array(ids))
	return err
}
