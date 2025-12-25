package networks

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
	GetByCode(ctx context.Context, code string) (*entity.Network, error)
	Create(ctx context.Context, network *entity.Network) error
	Update(ctx context.Context, network *entity.Network) error
	Delete(ctx context.Context, networks ...*entity.Network) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db}
}

func (r *repository) GetByCode(ctx context.Context, code string) (*entity.Network, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getNetworkByCode, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrNetworkNotFound
		}
		return nil, err
	}
	return m.convert(), nil
}

func (r *repository) Create(ctx context.Context, network *entity.Network) error {
	m := convert(network)
	_, err := r.db.Exec(ctx, createNetwork, m.Code, m.Name, m.Description, m.UpdatedBy)
	return err
}

func (r *repository) Update(ctx context.Context, network *entity.Network) error {
	m := convert(network)
	_, err := r.db.Exec(ctx, updateNetwork, m.Code, m.Name, m.Description, m.UpdatedBy)
	return err
}

func (r *repository) Delete(ctx context.Context, networks ...*entity.Network) error {
	codes := lo.Map(networks, func(item *entity.Network, _ int) string {
		return item.Code
	})
	_, err := r.db.Exec(ctx, deleteNetwork, pq.Array(codes))
	return err
}
