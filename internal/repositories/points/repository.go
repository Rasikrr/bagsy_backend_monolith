package points

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database"
	"github.com/Rasikrr/core/log"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
	"github.com/samber/lo"
)

type Repository interface {
	Create(ctx context.Context, point *entity.Point) error
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
	Update(ctx context.Context, point *entity.Point) error
	Delete(ctx context.Context, points ...*entity.Point) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, point *entity.Point) error {
	m, err := convert(point)
	if err != nil {
		return err
	}

	log.Infof(ctx, "Creating point with code: %+v", m)

	_, err = r.db.Exec(ctx, createPoint,
		m.Code,
		m.Name,
		m.Description,
		m.NetworkCode,
		m.CategoryID,
		string(m.Address),
		m.City,
		m.Active,
		string(m.Schedule),
		m.UpdatedBy,
	)
	return err
}

func (r *repository) Update(ctx context.Context, point *entity.Point) error {
	m, err := convert(point)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, updatePoint,
		m.Code,
		m.Name,
		m.Description,
		m.NetworkCode,
		m.CategoryID,
		string(m.Address),
		m.City,
		m.Active,
		string(m.Schedule),
		m.UpdatedBy,
	)
	return err
}

func (r *repository) GetByCode(ctx context.Context, code string) (*entity.Point, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getPointByCode, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrPointNotFound.WithError(err)
		}
		return nil, err
	}
	p := m.convert()
	return p, nil
}

func (r *repository) Delete(ctx context.Context, points ...*entity.Point) error {
	codes := lo.Map(points, func(item *entity.Point, _ int) string {
		return item.Code
	})
	_, err := r.db.Exec(ctx, deletePoint, pq.Array(codes))
	return err
}
