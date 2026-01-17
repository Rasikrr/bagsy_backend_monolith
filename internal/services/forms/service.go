package forms

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/form"
)

type formsRepository interface {
	Create(ctx context.Context, form *form.Form) error
}

type Service struct {
	formsRepo formsRepository
}

func NewService(formsRepo formsRepository) *Service {
	return &Service{formsRepo: formsRepo}
}

func (s *Service) Create(ctx context.Context, form *form.Form) error {
	return s.formsRepo.Create(ctx, form)
}
