package bagsies

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Rasikrr/core/log"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/core/database"
	"github.com/georgysavva/scany/v2/pgxscan"
)

type Repository interface {
	Create(ctx context.Context, b *entity.Bagsy) error
	GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, b *entity.Bagsy) error {
	m := convertToModel(b)

	log.Info(ctx, "create repo",
		log.Any("model", m),
	)

	_, err := r.db.Exec(
		ctx,
		createBagsy,
		m.ID,
		m.PointCode,
		m.ProviderPhone,
		m.UserPhone,
		m.FirstName,
		m.LastName,
		m.Description,
		m.Service,
		m.StartAt,
		m.EndAt,
		m.CreatedAt,
		m.UpdatedAt,
		m.UpdatedBy,
	)
	return err
}

func (r *repository) GetByParams(ctx context.Context, params *entity.BagsyParams) ([]*entity.Bagsy, error) {
	var ms models
	err := pgxscan.Select(
		ctx,
		r.db,
		&ms,
		getBagsyByParams,
		params.PointCode,
		params.StartAt,
		params.EndAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return ms.convert(), nil
}

func (r *repository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, deleteBagsy, id)
	return err
}
