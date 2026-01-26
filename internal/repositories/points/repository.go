package points

import (
	"context"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
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

func (r *Repository) Create(ctx context.Context, p *point.Point) error {
	m, err := convert(p)
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
			return point.ErrPointAlreadyExists
		}
		return domainErr.NewInternalError("failed to create point in db", err)
	}
	return nil
}

func (r *Repository) GetByCode(ctx context.Context, code string) (*point.Point, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getPointByCode, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, point.ErrPointNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get point from db", err)
	}
	p := m.convert()
	return p, nil
}

func (r *Repository) GetByNetworkCode(ctx context.Context, networkCode string) ([]*point.Point, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByNetworkCode, networkCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
	}
	return mm.convert(), nil
}

func (r *Repository) GetByCodes(ctx context.Context, codes []string) ([]*point.Point, error) {
	if len(codes) == 0 {
		return []*point.Point{}, nil
	}

	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getByCodes, pq.Array(codes))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*point.Point{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get points from db", err)
	}
	return mm.convert(), nil
}

func (r *Repository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var out bool
	err := pgxscan.Get(ctx, r.db, &out, existByCode, code)
	if err != nil {
		return false, domainErr.NewInternalError("failed to check if point exists by code", err)
	}
	return out, nil
}

func (r *Repository) Update(ctx context.Context, point *point.Point) error {
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

func (r *Repository) Delete(ctx context.Context, points ...*point.Point) error {
	codes := lo.Map(points, func(item *point.Point, _ int) string {
		return item.Code
	})
	_, err := r.db.Exec(ctx, deletePoint, pq.Array(codes))
	if err != nil {
		return domainErr.NewInternalError("failed to delete points from db", err)
	}
	return nil
}
