package networks

import (
	"context"
	"errors"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/networks"
	"github.com/jackc/pgx/v5"
)

type Service interface {
	GetByCode(ctx context.Context, code string) (*entity.Network, error)
}

type service struct {
	networksRepo networks.Repository
}

func NewService(repo networks.Repository) Service {
	return &service{networksRepo: repo}
}

func (s *service) GetByCode(ctx context.Context, code string) (*entity.Network, error) {
	network, err := s.networksRepo.GetByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNetworkNotFound
		}
		return nil, err
	}
	return network, nil
}
