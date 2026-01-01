package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
)

type networksRepository interface {
	GetByCode(ctx context.Context, code string) (*entity.Network, error)
}

type Service struct {
	networksRepo networksRepository
}

func NewService(networksRepo networksRepository) *Service {
	return &Service{
		networksRepo: networksRepo,
	}
}

func (s *Service) GetByCode(ctx context.Context, code string) (*entity.Network, error) {
	return s.networksRepo.GetByCode(ctx, code)
}
