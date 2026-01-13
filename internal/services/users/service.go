package users

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
)

type usersRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByParams(ctx context.Context, filter *query.UserFilter) ([]*entity.User, error)
	CountByFilter(ctx context.Context, filter *query.UserFilter) (int, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	Update(ctx context.Context, user *entity.User) error
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*entity.Point, error)
}

type mediaService interface {
	SetUserAvatar(ctx context.Context, phone string, mediaID uuid.UUID) error
	RemoveUserAvatar(ctx context.Context, phone string) error
	GenerateDownloadURL(ctx context.Context, fileKey string) (string, error)
}

type Service struct {
	txManager     database.TXManager
	usersRepo     usersRepository
	pointsService pointsService
	mediaService  mediaService
}

func NewService(
	txManager database.TXManager,
	usersRepo usersRepository,
	pointsService pointsService,
	mediaService mediaService,
) *Service {
	return &Service{
		txManager:     txManager,
		usersRepo:     usersRepo,
		pointsService: pointsService,
		mediaService:  mediaService,
	}
}

// Create создает нового пользователя
// Проверяет что пользователь с таким номером не существует или не активен
func (s *Service) Create(ctx context.Context, user *entity.User, password string) error {
	exists, err := s.usersRepo.ExistsByPhone(ctx, user.Phone)
	if err != nil {
		return err
	}

	if exists {
		return domainErr.NewConflictError("active user with this phone already exists", nil).
			WithDetail("phone", user.Phone)
	}

	if password != "" {
		passwordHash, hashErr := hash.Password(password)
		if hashErr != nil {
			return domainErr.NewInternalError("failed to hash password", hashErr)
		}
		user.Password = passwordHash
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

func (s *Service) GetUserProfile(ctx context.Context) (*dto.UserWithAvatar, error) {
	ses, err := session.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.GetByPhone(ctx, ses.Phone())
	if err != nil {
		return nil, err
	}

	userDTO, err := s.enrichUserWithAvatar(ctx, user)
	if err != nil {
		return nil, err
	}

	return userDTO, nil
}

func (s *Service) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return s.usersRepo.ExistsByPhone(ctx, phone)
}

// GetListByFilter возвращает список пользователей с пагинацией и учетом прав доступа
// Применяет ограничения на основе роли текущего пользователя
func (s *Service) GetListByFilter(ctx context.Context, requestedFilter *query.UserFilter) (*dto.PaginatedUsers, error) {
	userSession, err := session.GetSession(ctx)
	if err != nil {
		return nil, domainErr.NewUnauthorizedError("user session not found").WithError(err)
	}

	// Применяем ограничения на основе роли
	authorizedFilter, err := s.authorizeFilter(ctx, userSession, requestedFilter)
	if err != nil {
		return nil, err
	}

	// Получаем пользователей с пагинацией
	users, err := s.usersRepo.GetByParams(ctx, authorizedFilter)
	if err != nil {
		return nil, err
	}

	userDTOs, err := s.enrichUsersWithAvatars(ctx, users)
	if err != nil {
		return nil, err
	}

	// Получаем общее количество пользователей по фильтру (без limit/offset)
	total, err := s.usersRepo.CountByFilter(ctx, authorizedFilter)
	if err != nil {
		return nil, err
	}

	return &dto.PaginatedUsers{
		Users: userDTOs,
		Total: total,
	}, nil
}

func (s *Service) UpdateSchedule(ctx context.Context, phone string, schedule []entity.StaffSchedule) error {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}
	user.Schedule = schedule

	return s.usersRepo.Update(ctx, user)
}

func (s *Service) UpdateProfile(ctx context.Context, cmd *command.UpdateUserCommand) (*dto.UserWithAvatar, error) {
	ses, err := session.GetSession(ctx)
	if err != nil {
		return nil, err
	}

	var (
		updatedUser *entity.User
	)

	err = s.txManager.Transaction(ctx, database.TXOptions{
		IsolationLevel: coreEnum.IsoLevelReadCommited,
	},
		func(ctx context.Context) error {
			user, userErr := s.usersRepo.GetByPhone(ctx, ses.Phone())
			if userErr != nil {
				return userErr
			}

			user.Name = cmd.Name
			user.Surname = cmd.Surname

			if cmd.AvatarID != nil {
				err = s.mediaService.SetUserAvatar(ctx, user.Phone, *cmd.AvatarID)
				if err != nil {
					return err
				}
			}

			// 9. Обновить пользователя (name, surname)
			if err = s.usersRepo.Update(ctx, user); err != nil {
				return err
			}

			updatedUser = user
			return nil
		})

	if err != nil {
		return nil, err
	}

	// 10. После транзакции - повторно запросить пользователя чтобы получить file_key из JOIN
	updatedUser, err = s.usersRepo.GetByPhone(ctx, ses.Phone())
	if err != nil {
		return nil, err
	}

	userDTO, err := s.enrichUserWithAvatar(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	return userDTO, nil
}

func (s *Service) UpdateWithPassword(ctx context.Context, user *entity.User, rawPassword string) error {
	if rawPassword != "" {
		passwordHash, hashErr := hash.Password(rawPassword)
		if hashErr != nil {
			return domainErr.NewInternalError("failed to hash password", hashErr)
		}
		user.Password = passwordHash
	}
	return s.usersRepo.Update(ctx, user)
}

func (s *Service) UpdatePasswordByPhone(ctx context.Context, phone string, rawPassword string) error {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}
	passwordHash, hashErr := hash.Password(rawPassword)
	if hashErr != nil {
		return domainErr.NewInternalError("failed to hash password", hashErr)
	}
	user.Password = passwordHash
	return s.usersRepo.Update(ctx, user)
}

func (s *Service) RemoveAvatar(ctx context.Context) error {
	ses, err := session.GetSession(ctx)
	if err != nil {
		return err
	}
	return s.mediaService.RemoveUserAvatar(ctx, ses.Phone())
}

func (s *Service) enrichUserWithAvatar(ctx context.Context, user *entity.User) (*dto.UserWithAvatar, error) {
	if user == nil {
		return nil, domainErr.ErrUserNotFound
	}
	userDTO := &dto.UserWithAvatar{
		User: user,
	}
	if user.AvatarFileKey == nil {
		return userDTO, nil
	}
	url, err := s.mediaService.GenerateDownloadURL(ctx, *user.AvatarFileKey)
	if err != nil {
		return nil, err
	}
	userDTO.AvatarURL = &url
	return userDTO, nil
}

func (s *Service) enrichUsersWithAvatars(ctx context.Context, users []*entity.User) ([]*dto.UserWithAvatar, error) {
	if len(users) == 0 {
		return nil, nil
	}
	out := make([]*dto.UserWithAvatar, len(users))
	for i, user := range users {
		userDTO, err := s.enrichUserWithAvatar(ctx, user)
		if err != nil {
			return nil, err
		}
		out[i] = userDTO
	}
	return out, nil
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
