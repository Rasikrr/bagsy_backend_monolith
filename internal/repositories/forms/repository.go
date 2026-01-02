package forms

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/database/postgres"
)

type Repository struct {
	db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, form *entity.Form) error {
	m := toModel(form)
	_, err := r.db.Exec(ctx, insertForm, m.FirstName, m.LastName, m.Phone, m.Description, m.Role)
	if err != nil {
		return domainErr.NewInternalError("failed to create client form in db", err)
	}
	return nil
}
