package services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getServiceByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainErr.ErrServiceNotFound
		}
		return nil, err
	}
	return m.convert(), nil
}

func (r *Repository) Create(ctx context.Context, service *entity.Service) error {
	m := convert(service)
	err := r.db.QueryRow(ctx, createService,
		m.PointCode, m.CategoryID, m.SubcategoryID, m.Name,
		m.Description, m.DurationMinutes, m.Active, m.UpdatedBy,
	).Scan(&service.ID)
	return err
}

func (r *Repository) Update(ctx context.Context, service *entity.Service) error {
	m := convert(service)
	_, err := r.db.Exec(ctx, updateService,
		m.ID, m.PointCode, m.CategoryID, m.SubcategoryID, m.Name,
		m.Description, m.DurationMinutes, m.Active, m.UpdatedBy,
	)
	return err
}

func (r *Repository) Delete(ctx context.Context, services ...*entity.Service) error {
	ids := lo.Map(services, func(item *entity.Service, _ int) uuid.UUID {
		return item.ID
	})
	_, err := r.db.Exec(ctx, deleteService, pq.Array(ids))
	return err
}
