package masterservices

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/google/uuid"
)

type masterServicesRepository interface {
	GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*entity.MasterService, error)
}

type Service struct {
	masterServicesRepo masterServicesRepository
}

func NewService(repository masterServicesRepository) *Service {
	return &Service{
		masterServicesRepo: repository,
	}
}

func (s *Service) GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*entity.MasterService, error) {
	return s.masterServicesRepo.GetByMasterPhoneAndServiceID(ctx, phone, serviceID)
}
