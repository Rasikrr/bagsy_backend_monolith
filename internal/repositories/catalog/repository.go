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

func (r *Repository) SaveService(ctx context.Context, s *catalog.Service) error {
	m := fromServiceDomain(s)
	_, err := r.db.Exec(ctx, saveService,
		m.ID, m.LocationID, m.CategoryID, m.Name, m.Description, m.DurationMinutes,
		m.Color, m.SortOrder, m.Active, m.CreatedAt, m.UpdatedAt, m.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("save service: %w", err)
	}
	return nil
}

func (r *Repository) SaveEmployeeService(ctx context.Context, es *catalog.EmployeeService) error {
	m := fromEmployeeServiceDomain(es)
	_, err := r.db.Exec(ctx, saveEmployeeService,
		m.ID, m.EmployeeID, m.ServiceID, m.Price, m.Active, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("save employee service: %w", err)
	}
	return nil
}

func (r *Repository) GetServiceCategoryByID(ctx context.Context, id uuid.UUID) (*catalog.ServiceCategory, error) {
	var m serviceCategoryModel
	if err := pgxscan.Get(ctx, r.db, &m, getServiceCategoryByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, catalog.ErrServiceCategoryNotFound
		}
		return nil, fmt.Errorf("get service category by id: %w", err)
	}
	return m.toDomain(), nil
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
