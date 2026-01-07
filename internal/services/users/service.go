package users

import (
	"context"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
)

type usersRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByParams(ctx context.Context, filter *query.UserFilter) ([]*entity.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	Update(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, patch *users.UserUpdatePatch) error
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
}

type Service struct {
	usersRepo     usersRepository
	pointsService pointsService
}

func NewService(
	usersRepo usersRepository,
	pointsService pointsService,
) *Service {
	return &Service{
		usersRepo:     usersRepo,
		pointsService: pointsService,
	}
}

// Create создает нового пользователя
// Проверяет что пользователь с таким номером не существует или не активен
func (s *Service) Create(ctx context.Context, user *entity.User) error {
	exists, err := s.usersRepo.ExistsByPhone(ctx, user.Phone)
	if err != nil {
		return err
	}

	if exists {
		return domainErr.NewConflictError("active user with this phone already exists", nil).
			WithDetail("phone", user.Phone)
	}

	return s.usersRepo.Create(ctx, user)
}

func (s *Service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetUserProfile(ctx context.Context) (*entity.User, error) {
	ses, err := session.GetSession(ctx)
	if err != nil {
		return nil, err
	}
	return s.GetByPhone(ctx, ses.Phone())
}

func (s *Service) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return s.usersRepo.ExistsByPhone(ctx, phone)
}

// GetByFilter возвращает список пользователей с пагинацией и учетом прав доступа
// Применяет ограничения на основе роли текущего пользователя
func (s *Service) GetByFilter(ctx context.Context, requestedFilter *query.UserFilter) ([]*entity.User, error) {
	userSession, err := session.GetSession(ctx)
	if err != nil {
		return nil, domainErr.NewUnauthorizedError("user session not found").WithError(err)
	}

	// Применяем ограничения на основе роли
	authorizedFilter, err := s.authorizeFilter(ctx, userSession, requestedFilter)
	if err != nil {
		return nil, err
	}

	users, err := s.usersRepo.GetByParams(ctx, authorizedFilter)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// authorizeFilter применяет ограничения доступа на основе роли пользователя
// Возвращает модифицированный фильтр или ошибку при недостаточных правах
func (s *Service) authorizeFilter(
	ctx context.Context,
	userSession *session.Session,
	requestedFilter *query.UserFilter,
) (*query.UserFilter, error) {
	switch userSession.Role() {
	case enum.RoleAdmin:
		return requestedFilter, nil

	case enum.RoleNetManager, enum.RoleSelfOwner:
		// Могут получать пользователей только своей сети
		userNetworkCode := userSession.NetworkCode()

		// Если пытаются запросить другую сеть - запрещаем
		if requestedFilter.NetworkCode != nil && *requestedFilter.NetworkCode != userNetworkCode {
			return nil, domainErr.NewForbiddenError("cannot access users from other network").
				WithDetail("requested_network", *requestedFilter.NetworkCode).
				WithDetail("user_network", userNetworkCode)
		}

		// Принудительно устанавливаем свою сеть
		requestedFilter.NetworkCode = &userNetworkCode

		if requestedFilter.PointCode != nil {
			point, err := s.pointsService.GetByCode(ctx, *requestedFilter.PointCode)
			if err != nil {
				return nil, err
			}
			if point.NetworkCode != userNetworkCode {
				return nil, domainErr.NewForbiddenError("cannot access users from other network").
					WithDetail("point", *requestedFilter.PointCode).
					WithDetail("user_network", userNetworkCode)
			}
		}
		return requestedFilter, nil

	case enum.RoleManager:
		// Может получать пользователей только своей точки
		userPointCode := userSession.PointCode()

		// Если пытаются запросить другую точку - запрещаем
		if requestedFilter.PointCode != nil && *requestedFilter.PointCode != userPointCode {
			return nil, domainErr.NewForbiddenError("cannot access users from other point").
				WithDetail("requested_point", *requestedFilter.PointCode).
				WithDetail("user_point", userPointCode)
		}

		// Принудительно устанавливаем свою точку и сеть
		requestedFilter.PointCode = &userPointCode
		requestedFilter.NetworkCode = ptr.Pointer(userSession.NetworkCode())
		return requestedFilter, nil

	case enum.RoleStaff, enum.RoleUser:
		return nil, domainErr.NewForbiddenError("insufficient permissions to list users").
			WithDetail("role", userSession.Role().String())

	default:
		return nil, domainErr.NewForbiddenError("unknown role").
			WithDetail("role", userSession.Role().String())
	}
}

func (s *Service) Update(ctx context.Context, newUser *entity.User) error {
	exists, err := s.usersRepo.ExistsByPhone(ctx, newUser.Phone)
	if err != nil {
		return err
	}
	if !exists {
		return domainErr.ErrUserNotFound
	}

	oldUser, err := s.usersRepo.GetByPhone(ctx, newUser.Phone)
	if err != nil {
		return err
	}

	return s.usersRepo.Update(ctx, convertToUpdatedUser(oldUser, newUser))
}
