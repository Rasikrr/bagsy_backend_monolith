package services

import (
	"context"

	masterservice "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
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

type masterServicesRepository interface {
	GetByPointCodeAndServiceIDs(ctx context.Context, pointCode string, serviceID ...uuid.UUID) ([]*masterservice.MasterService, error)
}

type Service struct {
	masterServicesRepo masterServicesRepository
	serviceRepository  servicesRepository
}

func NewService(
	repository servicesRepository,
	masterServicesRepo masterServicesRepository,
) *Service {
	return &Service{
		masterServicesRepo: masterServicesRepo,
		serviceRepository:  repository,
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
	services, err := s.serviceRepository.GetByPointCode(ctx, pointCode)
	if err != nil {
		return nil, err
	}

	serviceIDs := lo.Map(services, func(item *service.Service, _ int) uuid.UUID { return item.ID })

	masterServices, err := s.masterServicesRepo.GetByPointCodeAndServiceIDs(ctx, pointCode, serviceIDs...)
	if err != nil {
		return nil, err
	}

	enrichWithPrices(services, masterServices)

	return services, nil
}

func enrichWithPrices(services []*service.Service, masterServices []*masterservice.MasterService) {
	pricesByServiceID := lo.GroupBy(masterServices, func(ms *masterservice.MasterService) uuid.UUID {
		return ms.ServiceID
	})

	for _, serv := range services {
		prices, ok := pricesByServiceID[serv.ID]
		if !ok || len(prices) == 0 {
			continue
		}

		minPrice := prices[0].Price
		maxPrice := prices[0].Price
		for _, p := range prices[1:] {
			if p.Price.LessThan(minPrice) {
				minPrice = p.Price
			}
			if p.Price.GreaterThan(maxPrice) {
				maxPrice = p.Price
			}
		}

		serv.MinPrice = &minPrice
		serv.MaxPrice = &maxPrice
	}
}

func (s *Service) Update(ctx context.Context, cmd *service.UpdateServiceCommand) error {
	return s.serviceRepository.Update(ctx, cmd)
}

func (s *Service) Delete(ctx context.Context, ids ...uuid.UUID) error {
	return s.serviceRepository.Delete(ctx, ids...)
}
