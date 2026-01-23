package services

import (
	"context"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/cockroachdb/errors"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error) {
	var m model
	err := pgxscan.Get(ctx, r.db, &m, getServiceByID, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, service.ErrServiceNotFound.WithError(err)
		}
		return nil, domainErr.NewInternalError("failed to get service from db", err)
	}
	return m.convert(), nil
}

func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*service.Service, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getServicesByIDs, pq.Array(ids))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, domainErr.NewInternalError("failed to get services from db", err)
	}
	return mm.convert(), nil
}

func (r *Repository) Create(ctx context.Context, service *service.Service) (uuid.UUID, error) {
	m := convert(service)
	var newID uuid.UUID
	err := pgxscan.Get(ctx, r.db, &newID, createService,
		m.PointCode, m.CategoryID, m.SubcategoryID, m.Name,
		m.Description, m.DurationMinutes, m.Color, m.Active, m.UpdatedBy,
	)
	if err != nil {
		return uuid.Nil, domainErr.NewInternalError("failed to create service in db", err)
	}
	return newID, nil
}

func (r *Repository) Update(ctx context.Context, service *service.Service) error {
	m := convert(service)
	_, err := r.db.Exec(ctx, updateService,
		m.ID, m.PointCode, m.CategoryID, m.SubcategoryID, m.Name,
		m.Description, m.DurationMinutes, m.Active, m.Color, m.UpdatedBy,
	)
	if err != nil {
		return domainErr.NewInternalError("failed to update service in db", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, ids ...uuid.UUID) error {
	_, err := r.db.Exec(ctx, deleteService, pq.Array(ids))
	if err != nil {
		return domainErr.NewInternalError("failed to delete services from db", err)
	}
	return nil
}

func (r *Repository) GetByPointCode(ctx context.Context, pointCode string) ([]*service.Service, error) {
	var mm models
	err := pgxscan.Select(ctx, r.db, &mm, getServicesByPointCode, pointCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []*service.Service{}, nil
		}
		return nil, domainErr.NewInternalError("failed to get services from db", err)
	}
	return mm.convert(), nil
}
