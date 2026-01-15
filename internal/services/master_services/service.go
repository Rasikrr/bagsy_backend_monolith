package masterservices

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/google/uuid"
)

type masterServicesRepository interface {
	GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*masterservice.MasterService, error)
}

type Service struct {
	masterServicesRepo masterServicesRepository
}

func NewService(repository masterServicesRepository) *Service {
	return &Service{
		masterServicesRepo: repository,
	}
}

func (s *Service) GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*masterservice.MasterService, error) {
	return s.masterServicesRepo.GetByMasterPhoneAndServiceID(ctx, phone, serviceID)
}
