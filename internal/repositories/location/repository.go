package location

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/location"
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
