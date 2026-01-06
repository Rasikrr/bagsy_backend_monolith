package networks

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/util/slug"
)

type networksRepository interface {
	Create(ctx context.Context, network *entity.Network) error
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

func (s *Service) Create(ctx context.Context, req *command.CreateNetworkCommand) error {
	ses, err := session.GetSession(ctx)
	if err != nil {
		return err
	}

	networkCode := slug.Generate(req.Name)
	network := &entity.Network{
		Code:        networkCode,
		Name:        req.Name,
		Description: &req.Description,
		CreatedBy:   ses.Phone(),
	}

	return s.networksRepo.Create(ctx, network)
}

func (s *Service) GetByCode(ctx context.Context, code string) (*entity.Network, error) {
	return s.networksRepo.GetByCode(ctx, code)
}
