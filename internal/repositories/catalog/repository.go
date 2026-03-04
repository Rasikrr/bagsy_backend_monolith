package catalog

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/catalog"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{
		db: db,
	}
}

// ─────────────────────────────────────────────────────────────────
// Service Repository Implementation
// ─────────────────────────────────────────────────────────────────

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*catalog.Service, error) {
	var m serviceModel
	if err := pgxscan.Get(ctx, r.db, &m, getServiceByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, catalog.ErrServiceNotFound
		}
		return nil, fmt.Errorf("get service by id: %w", err)
	}
	return m.toDomain()
}

// ─────────────────────────────────────────────────────────────────
// EmployeeService Repository Implementation
// ─────────────────────────────────────────────────────────────────

func (r *Repository) GetActiveByEmployeeAndService(ctx context.Context, employeeID, serviceID uuid.UUID) (*catalog.EmployeeService, error) {
	var m employeeServiceModel
	if err := pgxscan.Get(ctx, r.db, &m, getEmployeeServiceByEmployeeAndService, employeeID, serviceID); err != nil {
		if pgxscan.NotFound(err) {
			return nil, catalog.ErrEmployeeServiceNotFound
		}
		return nil, fmt.Errorf("get employee service: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) GetActiveByLocationAndService(ctx context.Context, locationID, serviceID uuid.UUID) ([]*catalog.EmployeeService, error) {
	var models []employeeServiceModel
	if err := pgxscan.Select(ctx, r.db, &models, getEmployeeServicesByLocationAndService, locationID, serviceID); err != nil {
		return nil, fmt.Errorf("select employee services: %w", err)
	}

	result := make([]*catalog.EmployeeService, 0, len(models))
	for _, m := range models {
		es, err := m.toDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, es)
	}
	return result, nil
}
