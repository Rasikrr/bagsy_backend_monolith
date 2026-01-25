package masterservices

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/master_service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/google/uuid"
)

const defaultActive = false

type masterServicesRepository interface {
	GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*masterservice.MasterService, error)
	GetByPointCodeAndServiceIDs(ctx context.Context, pointCode string, serviceID ...uuid.UUID) ([]*masterservice.MasterService, error)
	Create(ctx context.Context, masterService *masterservice.MasterService) error
}

type usersRepository interface {
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
}

type servicesRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*service.Service, error)
}

type Service struct {
	masterServicesRepo masterServicesRepository
	usersRepo          usersRepository
	servicesRepo       servicesRepository
}

func NewService(
	repository masterServicesRepository,
	usersRepo usersRepository,
	servicesRepo servicesRepository,
) *Service {
	return &Service{
		masterServicesRepo: repository,
		usersRepo:          usersRepo,
		servicesRepo:       servicesRepo,
	}
}

func (s *Service) Create(ctx context.Context, cmd *masterservice.CreateMasterServiceCommand) (*masterservice.MasterService, error) {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}

	targetPhone := act.Phone()
	if cmd.MasterPhone != nil {
		targetPhone = *cmd.MasterPhone
	}

	target, err := s.authorizeAndGetTarget(ctx, act, targetPhone)
	if err != nil {
		return nil, err
	}

	// Validate that the service exists at the target's point
	if target.PointCode == nil {
		return nil, domainErr.NewInvalidInputError("target user is not assigned to a point", nil)
	}
	svc, err := s.servicesRepo.GetByID(ctx, cmd.ServiceID)
	if err != nil {
		return nil, err
	}
	if svc.PointCode != *target.PointCode {
		return nil, domainErr.NewInvalidInputError("service does not belong to the target user's point", nil)
	}

	// Check for duplicate
	_, err = s.masterServicesRepo.GetByMasterPhoneAndServiceID(ctx, targetPhone, cmd.ServiceID)
	if err == nil {
		return nil, masterservice.ErrMasterServiceAlreadyExists
	}
	if !domainErr.IsNotFound(err) {
		return nil, err
	}

	var updatedBy *string
	actorPhone := act.Phone()
	if actorPhone != "" {
		updatedBy = &actorPhone
	}

	ms := &masterservice.MasterService{
		MasterPhone: targetPhone,
		ServiceID:   cmd.ServiceID,
		Price:       cmd.Price,
		Active:      defaultActive,
		UpdatedBy:   updatedBy,
	}

	if err = s.masterServicesRepo.Create(ctx, ms); err != nil {
		return nil, err
	}

	return ms, nil
}

//nolint
func (s *Service) authorizeAndGetTarget(ctx context.Context, act *actor.Actor, targetPhone string) (*user.User, error) {
	switch act.Role() {
	case user.RoleStaff:
		// Staff can only create for themselves
		if targetPhone != act.Phone() {
			return nil, domainErr.NewForbiddenError("staff can only create master service for themselves")
		}
		return s.usersRepo.GetByPhone(ctx, targetPhone)
	case user.RoleManager:
		// Manager cannot create for themselves, only for staff in their point
		if targetPhone == act.Phone() {
			return nil, domainErr.NewForbiddenError("manager cannot create master service for themselves")
		}
		target, err := s.usersRepo.GetByPhone(ctx, targetPhone)
		if err != nil {
			return nil, err
		}
		if target.PointCode == nil || *target.PointCode != act.PointCode() {
			return nil, domainErr.NewForbiddenError("target user is not in your point")
		}
		return target, nil
	case user.RoleNetManager:
		// NetManager/SelfOwner cannot create for themselves, only for staff in their network
		if targetPhone == act.Phone() {
			return nil, domainErr.NewForbiddenError("net manager cannot create master service for themselves")
		}
		target, err := s.usersRepo.GetByPhone(ctx, targetPhone)
		if err != nil {
			return nil, err
		}
		if target.NetworkCode == nil || *target.NetworkCode != act.NetworkCode() {
			return nil, domainErr.NewForbiddenError("target user is not in your network")
		}
		return target, nil
	case user.RoleSelfOwner:
		target, err := s.usersRepo.GetByPhone(ctx, targetPhone)
		if err != nil {
			return nil, err
		}
		if target.NetworkCode == nil || *target.NetworkCode != act.NetworkCode() {
			return nil, domainErr.NewForbiddenError("target user is not in your network")
		}
		return target, nil
	case user.RoleAdmin:
		// No restrictions, but still need the target user for point validation
		return s.usersRepo.GetByPhone(ctx, targetPhone)
	default:
		return nil, domainErr.NewForbiddenError("insufficient permissions")
	}
}

func (s *Service) GetByMasterPhoneAndServiceID(ctx context.Context, phone string, serviceID uuid.UUID) (*masterservice.MasterService, error) {
	return s.masterServicesRepo.GetByMasterPhoneAndServiceID(ctx, phone, serviceID)
}

func (s *Service) GetByPointCodeAndServiceID(ctx context.Context, pointCode string, serviceID uuid.UUID) ([]*masterservice.MasterService, error) {
	return s.masterServicesRepo.GetByPointCodeAndServiceIDs(ctx, pointCode, serviceID)
}
