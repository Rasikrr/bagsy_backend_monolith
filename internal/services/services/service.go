package services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

type servicesRepository interface {
	Create(ctx context.Context, cmd *service.CreateServiceCommand) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*service.Service, error)
	GetByPointCode(ctx context.Context, pointCode string) ([]*service.Service, error)
	Update(ctx context.Context, cmd *service.UpdateServiceCommand) error
	Delete(ctx context.Context, ids ...uuid.UUID) error
}

type Service struct {
	serviceRepository servicesRepository
}

func NewService(repository servicesRepository) *Service {
	return &Service{
		serviceRepository: repository,
	}
}

func (s *Service) Create(ctx context.Context, cmd *service.CreateServiceCommand) (uuid.UUID, error) {
	return s.serviceRepository.Create(ctx, cmd)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error) {
	return s.serviceRepository.GetByID(ctx, id)
}

func (s *Service) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*service.Service, error) {
	ids = lo.Uniq(ids)
	return s.serviceRepository.GetByIDs(ctx, ids)
}

func (s *Service) GetByPointCode(ctx context.Context, pointCode string) ([]*service.Service, error) {
	return s.serviceRepository.GetByPointCode(ctx, pointCode)
}

func (s *Service) Update(ctx context.Context, cmd *service.UpdateServiceCommand) error {
	return s.serviceRepository.Update(ctx, cmd)
}

func (s *Service) Delete(ctx context.Context, ids ...uuid.UUID) error {
	return s.serviceRepository.Delete(ctx, ids...)
}
