package points

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/Rasikrr/core/log"
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

func (r *Repository) Create(ctx context.Context, point *entity.Point) error {
	m, err := convert(point)
	if err != nil {
		return domainErr.NewInternalError("failed to convert point entity", err)
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
	if err != nil {
		if postgres.IsUniqueViolation(err) {
			return domainErr.ErrPointAlreadyExists
		}
		return domainErr.NewInternalError("failed to create point in db", err)
	}
	return nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*entity.Point, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getPointByCode, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrPointNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point from db", err)
	}
	p := m.convert()
	return p, nil
}

func (r *Repository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var out bool
	err := pgxscan.Get(ctx, r.db, &out, existByCode, code)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check if point exists by code", err)
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, point *entity.Point) error {
	m, err := convert(point)
	if err != nil {
		return domainErr.NewInternalError("failed to convert point entity", err)
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
	if err != nil {
		return domainErr.NewInternalError("failed to update point in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, points ...*entity.Point) error {
	codes := lo.Map(points, func(item *entity.Point, _ int) string {
		return item.Code
	})
	_, err := r.db.Exec(ctx, deletePoint, pq.Array(codes))
	if err != nil {
		return domainErr.NewInternalError("failed to delete points from db", err)
	}
	return nil
}
