package employee

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, identity.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("get employee by id: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) ExistsByPhone(ctx context.Context, phone shared.Phone) (bool, error) {
	var exists bool
	if err := pgxscan.Get(ctx, r.db, &exists, existsByPhone, phone.String()); err != nil {
		return false, fmt.Errorf("employee exists by phone: %w", err)
	}
	return exists, nil
}

func (r *Repository) GetByPhone(ctx context.Context, phone shared.Phone) (*identity.Employee, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByPhone, phone.String()); err != nil {
		if pgxscan.NotFound(err) {
			return nil, identity.ErrEmployeeNotFound
		}
		return nil, fmt.Errorf("get employee by phone: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) Save(ctx context.Context, emp *identity.Employee) error {
	m := fromDomain(emp)
	_, err := r.db.Exec(ctx, saveEmployee,
		m.ID,
		m.Phone,
		m.PasswordHash,
		m.FirstName,
		m.LastName,
		m.OrganizationID,
		m.LocationID,
		m.Role,
		m.CanProvideServices,
		m.CanManageLocationSchedule,
		m.Active,
		m.CreatedAt,
		m.UpdatedAt,
		m.DeletedAt,
		m.AvatarID,
	)
	if err != nil {
		return fmt.Errorf("save employee: %w", err)
	}
	return nil
}

func (r *Repository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int, error) {
	var count int
	if err := pgxscan.Get(ctx, r.db, &count, countByOrganization, orgID.String()); err != nil {
		return 0, fmt.Errorf("count employee by organization: %w", err)
	}
	return count, nil
}
