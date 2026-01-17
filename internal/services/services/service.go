package services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
)

type servicesRepository interface {
	Create(ctx context.Context, service *service.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error)
	GetByPointCode(ctx context.Context, pointCode string) ([]*service.Service, error)
	Update(ctx context.Context, service *service.Service) error
	Delete(ctx context.Context, service ...*service.Service) error
}

type Service struct {
	serviceRepository servicesRepository
}

func NewService(repository servicesRepository) *Service {
	return &Service{
		serviceRepository: repository,
	}
}

func (s *Service) Create(ctx context.Context, service *service.Service) error {
	return s.serviceRepository.Create(ctx, service)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error) {
	return s.serviceRepository.GetByID(ctx, id)
}

func (s *Service) GetByPointCode(ctx context.Context, pointCode string) ([]*service.Service, error) {
	return s.serviceRepository.GetByPointCode(ctx, pointCode)
}
