package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/core/database"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Repository interface {
	GetByCode(ctx context.Context, code string) (*entity.Network, error)
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
		return nil, err
	}
	return m.convert(), nil
}
