package location

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
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

func (r *Repository) Save(ctx context.Context, l *location.Location) error {
	m := fromDomain(l)
	_, err := r.db.Exec(ctx, saveLocation,
		m.ID, m.OrganizationID, m.CategoryID, m.Name, m.Description, m.Phone, m.Slug,
		m.City, m.Street, m.Building, m.Details,
		m.Longitude, m.Latitude, m.Active, m.ScheduleType, m.SlotDurationMinutes,
		m.CreatedAt, m.UpdatedAt, m.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("save location: %w", err)
	}
	return nil
}

func (r *Repository) CountByOrganization(ctx context.Context, organizationID uuid.UUID) (int, error) {
	var count int
	if err := pgxscan.Get(ctx, r.db, &count, countByOrganization, organizationID); err != nil {
		return 0, fmt.Errorf("count locations by organization: %w", err)
	}
	return count, nil
}

func (r *Repository) GetByFilter(ctx context.Context, filter *location.Filter) (*shared.Page[*location.Location], error) {
	base := buildLocationFilterBase(filter)

	countSQL, countArgs, err := base.Columns("COUNT(*)").ToSql()
	if err != nil {
		return nil, fmt.Errorf("build count query: %w", err)
	}

	var total int
	if err = pgxscan.Get(ctx, r.db, &total, countSQL, countArgs...); err != nil {
		return nil, fmt.Errorf("count locations by filter: %w", err)
	}

	dataSQL, dataArgs, err := base.
		Columns(locationColumns).
		OrderBy(filter.OrderBy.String() + " " + filter.SortOrder.String()).
		Limit(filter.Limit).
		Offset(filter.Offset).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("build select query: %w", err)
	}

	var models []model
	if err = pgxscan.Select(ctx, r.db, &models, dataSQL, dataArgs...); err != nil {
		return nil, fmt.Errorf("select locations by filter: %w", err)
	}

	items := make([]*location.Location, 0, len(models))
	for _, m := range models {
		var loc *location.Location
		loc, err = m.toDomain()
		if err != nil {
			return nil, err
		}
		items = append(items, loc)
	}

	return &shared.Page[*location.Location]{
		Items: items,
		Total: total,
	}, nil
}

const locationColumns = `id, organization_id, category_id, name, description, phone, slug,
	city, address_street, address_building, address_details,
	longitude, latitude, active, schedule_type, slot_duration_minutes,
	created_at, updated_at, deleted_at`

func buildLocationFilterBase(filter *location.Filter) sq.SelectBuilder {
	builder := sq.Select().
		PlaceholderFormat(sq.Dollar).
		From("locations").
		Where(sq.Eq{"organization_id": filter.OrganizationID}).
		Where("deleted_at IS NULL")

	if filter.Active != nil {
		builder = builder.Where(sq.Eq{"active": *filter.Active})
	}

	return builder
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*location.Location, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, location.ErrLocationNotFound
		}
		return nil, fmt.Errorf("get location by id: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) GetBySlug(ctx context.Context, slug string) (*location.Location, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getBySlug, slug); err != nil {
		if pgxscan.NotFound(err) {
			return nil, location.ErrLocationNotFound
		}
		return nil, fmt.Errorf("get location by slug: %w", err)
	}
	return m.toDomain()
}
