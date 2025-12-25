package forms

import (
	"context"

	"github.com/Rasikrr/core/database"
)

type Repository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateClient(ctx context.Context, firstName, lastName, phone, description string, role string) error {
	_, err := r.db.Exec(ctx, insertForm, firstName, lastName, phone, description, role)
	return err
}
