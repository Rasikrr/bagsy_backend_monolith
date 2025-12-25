package bagsies

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/samber/lo"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Bagsy, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrBagsyNotFound
		}
		return nil, err
	}
	return m.convert()
}

func (r *Repository) Create(ctx context.Context, b *entity.Bagsy) error {
	m := convertToModel(b)

	_, err := r.db.Exec(
		ctx,
		create,
		m.ID,
		m.PointCode,
		m.ClientPhone,
		m.Status,
		m.Price,
		m.MasterPhone,
		m.ServiceID,
		m.StartAt,
		m.EndAt,
		m.CreatedAt,
		m.UpdatedAt,
		m.UpdatedBy,
	)
	return err
}

func (r *Repository) Update(ctx context.Context, b *entity.Bagsy) error {
	m := convertToModel(b)
	_, err := r.db.Exec(
		ctx,
		update,
		m.ID,
		m.PointCode,
		m.ClientPhone,
		m.Status,
		m.Price,
		m.MasterPhone,
		m.ServiceID,
		m.StartAt,
		m.EndAt,
		m.CreatedAt,
		m.UpdatedAt,
		m.UpdatedBy,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, bagsies ...*entity.Bagsy) error {
	ids := lo.Map(bagsies, func(item *entity.Bagsy, _ int) uuid.UUID {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteByIDs, pq.Array(ids))
	return err
}
