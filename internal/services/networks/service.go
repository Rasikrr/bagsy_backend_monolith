package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
)

type networksRepository interface {
	Create(ctx context.Context, network *network.Network) error
	GetByCode(ctx context.Context, code string) (*network.Network, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
}

type Service struct {
	networksRepo networksRepository
}

func NewService(networksRepo networksRepository) *Service {
	return &Service{
		networksRepo: networksRepo,
	}
}

func (s *Service) RegisterNewNetwork(ctx context.Context, req *network.CreateNetworkCommand, createdBy string) (*network.Network, error) {
	networkCode := slug.Generate(req.Name)
	newNetwork := &network.Network{
		Code:        networkCode,
		Name:        req.Name,
		Description: &req.Description,
		CreatedBy:   createdBy,
	}

	err := s.networksRepo.Create(ctx, newNetwork)
	if err != nil {
		return nil, err
	}
	return newNetwork, nil
}

func (s *Service) GetByCode(ctx context.Context, code string) (*network.Network, error) {
	return s.networksRepo.GetByCode(ctx, code)
}

func (s *Service) ExistsByCode(ctx context.Context, code string) (bool, error) {
	return s.networksRepo.ExistsByCode(ctx, code)
}
