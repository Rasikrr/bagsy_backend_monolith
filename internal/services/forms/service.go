package forms

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type formsRepository interface {
	Create(ctx context.Context, form *entity.Form) error
}

type Service struct {
	formsRepo formsRepository
}

func NewService(formsRepo formsRepository) *Service {
	return &Service{formsRepo: formsRepo}
}

func (s *Service) Create(ctx context.Context, form *entity.Form) error {
	return s.formsRepo.Create(ctx, form)
}
