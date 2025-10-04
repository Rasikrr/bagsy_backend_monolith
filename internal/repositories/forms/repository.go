package forms

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	"github.com/Rasikrr/core/database"
)

type Repository interface {
	CreateClient(ctx context.Context, firstName, lastName, phone, description string, role enum.Role) error
}

type repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) Repository {
	return &repository{db: db}
}

func (r *repository) CreateClient(ctx context.Context, firstName, lastName, phone, description string, role enum.Role) error {
	_, err := r.db.Exec(ctx, insertForm, firstName, lastName, phone, description, role.String())
	return err
}
