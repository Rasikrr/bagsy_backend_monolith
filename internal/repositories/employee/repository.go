package employee

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/database/postgres"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/lib/pq"
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

func (r *Repository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*identity.Employee, error) {
	var models []model
	if err := pgxscan.Select(ctx, r.db, &models, getByIDs, pq.Array(ids)); err != nil {
		return nil, fmt.Errorf("select employees by ids: %w", err)
	}

	result := make([]*identity.Employee, 0, len(models))
	for _, m := range models {
		emp, err := m.toDomain()
		if err != nil {
			return nil, err
		}
		result = append(result, emp)
	}
	return result, nil
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

func (r *Repository) GetByFilter(ctx context.Context, filter *identity.EmployeeFilter) (*identity.EmployeePage, error) {
	base := buildFilterBase(filter)

	countSQL, countArgs, err := base.Columns("COUNT(*)").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err = pgxscan.Get(ctx, r.db, &total, countSQL, countArgs...); err != nil {
		return nil, fmt.Errorf("count employees by filter: %w", err)
	}

	dataSQL, dataArgs, err := base.
		Columns(employeeColumns).
		OrderBy(filter.OrderBy.String() + " " + filter.SortOrder.String()).
		Limit(filter.Limit).
		Offset(filter.Offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	var models []model
	if err = pgxscan.Select(ctx, r.db, &models, dataSQL, dataArgs...); err != nil {
		return nil, fmt.Errorf("select employees by filter: %w", err)
	}

	items := make([]*identity.Employee, 0, len(models))
	for _, m := range models {
		var emp *identity.Employee
		emp, err = m.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, emp)
	}

	return &identity.EmployeePage{
		Items: items,
		Total: total,
	}, nil
}

const employeeColumns = `id, phone, password_hash, first_name, last_name, avatar_id,
	organization_id, location_id, role,
	can_provide_services, can_manage_location_schedule,
	active, created_at, updated_at, deleted_at`

func buildFilterBase(filter *identity.EmployeeFilter) sq.SelectBuilder {
	builder := sq.Select().
		PlaceholderFormat(sq.Dollar).
		From("employees").
		Where(sq.Eq{"organization_id": filter.OrganizationID}).
		Where("deleted_at IS NULL")

	if filter.LocationID != nil {
		builder = builder.Where(sq.Eq{"location_id": *filter.LocationID})
	}

	if len(filter.Roles) > 0 {
		roles := make([]string, len(filter.Roles))
		for i, r := range filter.Roles {
			roles[i] = r.String()
		}
		builder = builder.Where(sq.Eq{"role": roles})
	}

	if filter.PhoneSearch != nil {
		builder = builder.Where("phone ILIKE '%' || ? || '%'", *filter.PhoneSearch)
	}

	if filter.Active != nil {
		builder = builder.Where(sq.Eq{"active": *filter.Active})
	}

	return builder
}
