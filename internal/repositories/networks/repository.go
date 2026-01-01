package networks

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

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Network, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getNetworkByCode, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrNetworkNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get network from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, network *entity.Network) error {
	m := convert(network)
	_, err := r.db.Exec(ctx, createNetwork, m.Code, m.Name, m.Description, m.UpdatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to create network in db", err)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, network *entity.Network) error {
	m := convert(network)
	_, err := r.db.Exec(ctx, updateNetwork, m.Code, m.Name, m.Description, m.UpdatedBy)
	if err != nil {
		return domainErr.NewInternalError("failed to update network in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, networks ...*entity.Network) error {
	codes := lo.Map(networks, func(item *entity.Network, _ int) string {
		return item.Code
	})
	_, err := r.db.Exec(ctx, deleteNetwork, pq.Array(codes))
	if err != nil {
		return domainErr.NewInternalError("failed to delete networks from db", err)
	}
	return nil
}
