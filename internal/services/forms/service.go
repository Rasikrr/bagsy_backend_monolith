package forms

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/forms"
)

type Service interface {
	CreateClient(ctx context.Context, firstName, lastName, phone, description string, role string) error
}

type service struct {
	formsRepo forms.Repository
}

func NewService(formsRepo forms.Repository) Service {
	return &service{formsRepo: formsRepo}
}

func (s *service) CreateClient(ctx context.Context, firstName, lastName, phone, description string, role string) error {
	return s.formsRepo.CreateClient(ctx, firstName, lastName, phone, description, role)
}
