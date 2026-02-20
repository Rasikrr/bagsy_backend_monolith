package customer

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

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*identity.Customer, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByID, id); err != nil {
		if pgxscan.NotFound(err) {
			return nil, identity.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("get customer by id: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) GetByPhone(ctx context.Context, phone shared.Phone) (*identity.Customer, error) {
	var m model
	if err := pgxscan.Get(ctx, r.db, &m, getByPhone, phone.String()); err != nil {
		if pgxscan.NotFound(err) {
			return nil, identity.ErrCustomerNotFound
		}
		return nil, fmt.Errorf("get customer by phone: %w", err)
	}
	return m.toDomain()
}

func (r *Repository) Save(ctx context.Context, c *identity.Customer) error {
	m := fromDomain(c)
	_, err := r.db.Exec(ctx, saveCustomer,
		m.ID, m.Phone, m.FirstName, m.LastName, m.BirthDate,
		m.CreatedAt, m.UpdatedAt, m.DeletedAt,
	)
	if err != nil {
		return fmt.Errorf("save customer: %w", err)
	}
	return nil
}
