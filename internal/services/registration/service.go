package registration

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/network"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	networksS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/networks"
	usersS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
)

type networkService interface {
	RegisterNewNetwork(ctx context.Context, in *network.CreateNetworkCommand, ownerPhone string) (*network.Network, error)
}

type usersService interface {
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
	CreateOwner(ctx context.Context, cmd *user.CreateOwnerCommand) (*user.User, error)
	PromoteToOwner(ctx context.Context, u *user.User, cmd *user.PromoteToOwnerCommand) (*user.User, error)
	CreateStaff(ctx context.Context, cmd *user.CreateStaffCommand) (*user.User, error)
	PromoteToStaff(ctx context.Context, u *user.User, cmd *user.PromoteToStaffCommand) (*user.User, error)
}

type Service struct {
	usersService   usersService
	networkService networkService
	txManager      database.TXManager
}

func NewService(
	txManager database.TXManager,
	users *usersS.Service,
	networks *networksS.Service,
) *Service {
	return &Service{
		usersService:   users,
		networkService: networks,
		txManager:      txManager,
	}
}

// RegisterNewOwner - регистрирует владельца и его сеть
func (s *Service) RegisterNewOwner(ctx context.Context, cmd *auth.RegisterManagementCommand) (*user.User, error) {
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	var (
		owner *user.User
	)

	err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		net, netErr := s.networkService.RegisterNewNetwork(
			txCtx,
			s.mapRegisterNetworkToCommand(cmd.NetworkRegisterInfo),
			cmd.Phone,
		)
		if netErr != nil {
			return netErr
		}
		var err error

		owner, err = s.usersService.GetByPhone(txCtx, cmd.Phone)
		if err != nil {
			if !domainErr.IsNotFound(err) {
				return err
			}
			// Надо создать
			createOwner := &user.CreateOwnerCommand{
				Name:        cmd.Name,
				Surname:     cmd.Surname,
				Password:    cmd.Password,
				Phone:       cmd.Phone,
				Role:        cmd.Role,
				NetworkCode: net.Code,
			}
			owner, err = s.usersService.CreateOwner(txCtx, createOwner)
		} else {
			owner, err = s.usersService.PromoteToOwner(txCtx, owner, &user.PromoteToOwnerCommand{
				Name:        cmd.Name,
				Surname:     cmd.Surname,
				Password:    cmd.Password,
				Role:        cmd.Role,
				NetworkCode: net.Code,
			})
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return owner, nil
}

// RegisterNewStaff - регистрирует стафф(manger/staff) и привязывает его к точке
func (s *Service) RegisterNewStaff(ctx context.Context, cmd *auth.RegisterStaffCommand, rawPassword string) (*user.User, error) {
	txOpts := database.TXOptions{IsolationLevel: coreEnum.IsoLevelReadCommited}

	var (
		staff *user.User
	)

	err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
		var err error
		staff, err = s.usersService.GetByPhone(txCtx, cmd.Phone)
		if err != nil {
			if !domainErr.IsNotFound(err) {
				return err
			}
			// Значит надо создать нового юзера
			staff, err = s.usersService.CreateStaff(txCtx, s.mapRegisterStaffToCreateStaff(cmd, rawPassword))
		} else {
			staff, err = s.usersService.PromoteToStaff(txCtx, staff, &user.PromoteToStaffCommand{
				Name:        cmd.Name,
				Surname:     cmd.Surname,
				Password:    rawPassword,
				Role:        cmd.Role,
				NetworkCode: cmd.NetworkCode,
				PointCode:   cmd.PointCode,
			})
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *Service) mapRegisterStaffToCreateStaff(cmd *auth.RegisterStaffCommand, rawPassword string) *user.CreateStaffCommand {
	return &user.CreateStaffCommand{
		Name:        cmd.Name,
		Surname:     cmd.Surname,
		Password:    rawPassword,
		Phone:       cmd.Phone,
		Role:        cmd.Role,
		NetworkCode: cmd.NetworkCode,
		PointCode:   cmd.PointCode,
	}
}

func (s *Service) mapRegisterNetworkToCommand(info *auth.RegisterNetworkInfo) *network.CreateNetworkCommand {
	return &network.CreateNetworkCommand{
		Name:        info.Name,
		Description: info.Description,
	}
}
