package services

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/google/uuid"
)

type servicesRepository interface {
	Create(ctx context.Context, service *entity.Service) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error)
	Update(ctx context.Context, service *entity.Service) error
	Delete(ctx context.Context, service ...*entity.Service) error
}

type Service struct {
	serviceRepository servicesRepository
}

func NewService(repository servicesRepository) *Service {
	return &Service{
		serviceRepository: repository,
	}
}

func (s *Service) Create(ctx context.Context, service *entity.Service) error {
	return s.serviceRepository.Create(ctx, service)
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*entity.Service, error) {
	return s.serviceRepository.GetByID(ctx, id)
}
