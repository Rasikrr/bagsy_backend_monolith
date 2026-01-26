package services

import (
	"context"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	masterservice "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

const defaultServiceActive = false

type servicesRepository interface {
	Create(ctx context.Context, service *service.Service) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*service.Service, error)
	GetByPointCode(ctx context.Context, pointCode string, isActive *bool) ([]*service.Service, error)
}

type masterServicesRepository interface {
	GetByPointCodeAndServiceIDs(ctx context.Context, pointCode string, serviceID ...uuid.UUID) ([]*masterservice.MasterService, error)
}

type serviceCategoriesRepository interface {
	GetByID(ctx context.Context, id int) (*service.Category, error)
}

type serviceSubcategoriesRepository interface {
	GetByID(ctx context.Context, id int) (*service.Subcategory, error)
}

type Service struct {
	masterServicesRepo       masterServicesRepository
	serviceRepository        servicesRepository
	serviceCategoriesRepo    serviceCategoriesRepository
	serviceSubcategoriesRepo serviceSubcategoriesRepository
}

func NewService(
	repository servicesRepository,
	masterServicesRepo masterServicesRepository,
	serviceCategoriesRepo serviceCategoriesRepository,
	serviceSubcategoriesRepo serviceSubcategoriesRepository,
) *Service {
	return &Service{
		masterServicesRepo:       masterServicesRepo,
		serviceRepository:        repository,
		serviceCategoriesRepo:    serviceCategoriesRepo,
		serviceSubcategoriesRepo: serviceSubcategoriesRepo,
	}
}

func (s *Service) Create(ctx context.Context, cmd *service.CreateServiceCommand) (uuid.UUID, error) {
	// Validate category exists
	_, err := s.serviceCategoriesRepo.GetByID(ctx, cmd.CategoryID)
	if err != nil {
		return uuid.Nil, err
	}

	// Validate subcategory exists and belongs to category (if provided)
	if cmd.SubcategoryID != nil {
		subcategory, subcatErr := s.serviceSubcategoriesRepo.GetByID(ctx, *cmd.SubcategoryID)
		if subcatErr != nil {
			return uuid.Nil, subcatErr
		}
		if subcategory.CategoryID != cmd.CategoryID {
			return uuid.Nil, domainErr.NewInvalidInputError("subcategory does not belong to category", nil)
		}
	}

	svc := &service.Service{
		PointCode:       cmd.PointCode,
		CategoryID:      cmd.CategoryID,
		SubcategoryID:   cmd.SubcategoryID,
		Name:            cmd.Name,
		Description:     cmd.Description,
		DurationMinutes: cmd.DurationMinutes,
		Active:          defaultServiceActive,
		UpdatedBy:       &cmd.UpdatedBy,
		Color:           cmd.Color,
	}
	id, err := s.serviceRepository.Create(ctx, svc)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error) {
	return s.serviceRepository.GetByID(ctx, id)
}
func (s *Service) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*service.Service, error) {
	ids = lo.Uniq(ids)
	return s.serviceRepository.GetByIDs(ctx, ids)
}

func (s *Service) GetByPointCode(ctx context.Context, pointCode string, isActive *bool) ([]*service.Service, error) {
	services, err := s.serviceRepository.GetByPointCode(ctx, pointCode, isActive)
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
