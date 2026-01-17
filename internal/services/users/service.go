package users

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/actor"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/user"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/hash"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/util/ptr"
	"github.com/Rasikrr/core/database"
	coreEnum "github.com/Rasikrr/core/enum"
	"github.com/google/uuid"
)

type usersRepository interface {
	Create(ctx context.Context, user *user.User) error
	GetByPhone(ctx context.Context, phone string) (*user.User, error)
	GetByPhones(ctx context.Context, phones []string) ([]*user.User, error)
	GetByParams(ctx context.Context, filter *user.Filter) ([]*user.User, error)
	CountByFilter(ctx context.Context, filter *user.Filter) (int, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	Update(ctx context.Context, user *user.User) error
}

type pointsService interface {
	GetByCode(ctx context.Context, code string) (*point.Point, error)
}

type userPhotosService interface {
	GetAvatarURL(ctx context.Context, fileKey string) (string, error)
	SetUserAvatar(ctx context.Context, phone string, mediaID uuid.UUID) error
	RemoveUserAvatar(ctx context.Context, phone string) error
}

type Service struct {
	txManager         database.TXManager
	usersRepo         usersRepository
	pointsService     pointsService
	userPhotosService userPhotosService
}

func NewService(
	txManager database.TXManager,
	usersRepo usersRepository,
	pointsService pointsService,
	mediaService userPhotosService,
) *Service {
	return &Service{
		txManager:         txManager,
		usersRepo:         usersRepo,
		pointsService:     pointsService,
		userPhotosService: mediaService,
	}
}

// CreateOwner создает нового пользователя (net_manager/self_owner)
// Должен вызываться в Registration service
func (s *Service) CreateOwner(ctx context.Context, cmd *user.CreateOwnerCommand) (*user.User, error) {
	if !cmd.Role.OneOf(user.RoleNetManager, user.RoleSelfOwner) {
		return nil, domainErr.NewInvalidInputError(
			"invalid role for owner registration",
			nil,
		).WithDetail("role", cmd.Role.String())
	}

	owner := &user.User{
		Name:        cmd.Name,
		Surname:     cmd.Surname,
		Role:        cmd.Role,
		NetworkCode: &cmd.NetworkCode,
		Active:      true,
	}
	passHash, err := hash.Password(cmd.Password)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to hash password", err)
	}
	owner.PasswordHash = passHash

	err = s.usersRepo.Create(ctx, owner)
	if err != nil {
		return nil, err
	}
	return owner, nil
}

func (s *Service) PromoteToOwner(
	ctx context.Context,
	u *user.User,
	cmd *user.PromoteToOwnerCommand,
) (*user.User, error) {
	return s.promote(ctx, u, cmd.ToPromoteNewLocation())
}

func (s *Service) PromoteToStaff(
	ctx context.Context,
	u *user.User,
	cmd *user.PromoteToStaffCommand,
) (*user.User, error) {
	return s.promote(ctx, u, cmd.ToPromoteNewLocation())
}

func (s *Service) promote(
	ctx context.Context,
	u *user.User,
	cmd *user.PromoteToNewLocationCommand,
) (*user.User, error) {
	err := cmd.Validate()
	if err != nil {
		return nil, err
	}
	if u.IsAssignedToLocation() {
		return nil, user.ErrUserBelongsToLocation
	}
	u.DetachFromLocation()

	u.Name = cmd.Name
	u.Surname = cmd.Surname
	u.Role = cmd.Role
	u.NetworkCode = &cmd.NetworkCode
	u.Active = true

	if cmd.PointCode != nil {
		u.PointCode = cmd.PointCode
	}

	passHash, err := hash.Password(cmd.Password)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to hash password", err)
	}
	u.PasswordHash = passHash

	err = s.usersRepo.Update(ctx, u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (s *Service) CreateStaff(ctx context.Context, cmd *user.CreateStaffCommand) (*user.User, error) {
	if !cmd.Role.OneOf(user.RoleManager, user.RoleStaff) {
		return nil, domainErr.NewInvalidInputError(
			"invalid role for staff registration",
			nil,
		).WithDetail("role", cmd.Role.String())
	}
	staff := &user.User{
		Name:        cmd.Name,
		Surname:     cmd.Surname,
		Role:        cmd.Role,
		NetworkCode: &cmd.NetworkCode,
		PointCode:   &cmd.PointCode,
		Active:      true,
	}
	passHash, err := hash.Password(cmd.Password)
	if err != nil {
		return nil, domainErr.NewInternalError("failed to hash password", err)
	}
	staff.PasswordHash = passHash

	err = s.usersRepo.Create(ctx, staff)
	if err != nil {
		return nil, err
	}
	return staff, nil
}

func (s *Service) CreateUser(ctx context.Context, cmd *user.CreateUserCommand) (*user.User, error) {
	newUser := &user.User{
		Name:    cmd.Name,
		Surname: cmd.Surname,
		Phone:   cmd.Phone,
		Active:  true,
	}

	err := s.usersRepo.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *Service) GetByPhone(ctx context.Context, phone string) (*user.User, error) {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) GetByPhones(ctx context.Context, phones []string) ([]*user.User, error) {
	return s.usersRepo.GetByPhones(ctx, phones)
}

func (s *Service) GetUserProfile(ctx context.Context) (*user.User, error) {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}

	user, err := s.usersRepo.GetByPhone(ctx, act.Phone())
	if err != nil {
		return nil, err
	}

	err = s.enrichUserWithAvatar(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return s.usersRepo.ExistsByPhone(ctx, phone)
}

// GetListByFilter возвращает список пользователей с пагинацией и учетом прав доступа
// Применяет ограничения на основе роли текущего пользователя
func (s *Service) GetListByFilter(ctx context.Context, requestedFilter *user.Filter) (*query.Page[*user.User], error) {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}

	// Применяем ограничения на основе роли
	authorizedFilter, err := s.authorizeFilter(ctx, act, requestedFilter)
	if err != nil {
		return nil, err
	}

	// Получаем пользователей с пагинацией
	users, err := s.usersRepo.GetByParams(ctx, authorizedFilter)
	if err != nil {
		return nil, err
	}

	err = s.enrichUsersWithAvatars(ctx, users)
	if err != nil {
		return nil, err
	}

	// Получаем общее количество пользователей по фильтру (без limit/offset)
	total, err := s.usersRepo.CountByFilter(ctx, authorizedFilter)
	if err != nil {
		return nil, err
	}

	return query.NewPage(users, total), nil
}

func (s *Service) UpdateSchedule(ctx context.Context, phone string, schedule user.Schedule) error {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return err
	}
	user.Schedule = schedule

	return s.usersRepo.Update(ctx, user)
}

func (s *Service) UpdateProfile(ctx context.Context, cmd *user.UpdateUserCommand) (*user.User, error) {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return nil, err
	}

	var (
		updatedUser *user.User
	)

	err = s.txManager.Transaction(ctx, database.TXOptions{
		IsolationLevel: coreEnum.IsoLevelReadCommited,
	},
		func(txCtx context.Context) error {
			user, userErr := s.usersRepo.GetByPhone(txCtx, act.Phone())
			if userErr != nil {
				return userErr
			}

			if cmd.Name != "" {
				user.Name = cmd.Name
			}
			if cmd.Surname != "" {
				user.Surname = cmd.Surname
			}

			if cmd.AvatarID != nil {
				err = s.userPhotosService.SetUserAvatar(txCtx, user.Phone, *cmd.AvatarID)
				if err != nil {
					return err
				}
			}

			// 9. Обновить пользователя (name, surname)
			err = s.usersRepo.Update(txCtx, user)
			if err != nil {
				return err
			}

			updatedUser = user
			return nil
		})

	if err != nil {
		return nil, err
	}

	// 10. После транзакции - повторно запросить пользователя чтобы получить file_key из JOIN
	updatedUser, err = s.usersRepo.GetByPhone(ctx, act.Phone())
	if err != nil {
		return nil, err
	}

	err = s.enrichUserWithAvatar(ctx, updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (s *Service) UpdateWithPassword(ctx context.Context, user *user.User, rawPassword string) error {
	if rawPassword != "" {
		passwordHash, hashErr := hash.Password(rawPassword)
		if hashErr != nil {
			return domainErr.NewInternalError("failed to hash password", hashErr)
		}
		user.PasswordHash = passwordHash
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
	user.PasswordHash = passwordHash
	return s.usersRepo.Update(ctx, user)
}

func (s *Service) RemoveAvatar(ctx context.Context) error {
	act, err := actor.GetActor(ctx)
	if err != nil {
		return err
	}
	return s.userPhotosService.RemoveUserAvatar(ctx, act.Phone())
}

func (s *Service) enrichUserWithAvatar(ctx context.Context, u *user.User) error {
	if u == nil {
		return user.ErrUserNotFound
	}
	if u.Avatar == nil || u.Avatar.FileKey == nil {
		return nil
	}

	url, err := s.userPhotosService.GetAvatarURL(ctx, ptr.Deref(u.Avatar.FileKey))
	if err != nil {
		return err
	}
	u.Avatar.URL = url
	return nil
}

func (s *Service) enrichUsersWithAvatars(ctx context.Context, users []*user.User) error {
	if len(users) == 0 {
		return nil
	}
	for _, u := range users {
		err := s.enrichUserWithAvatar(ctx, u)
		if err != nil {
			return err
		}
	}
	return nil
}

// authorizeFilter применяет ограничения доступа на основе роли пользователя
// Возвращает модифицированный фильтр или ошибку при недостаточных правах
func (s *Service) authorizeFilter(
	ctx context.Context,
	userSession *actor.Actor,
	requestedFilter *user.Filter,
) (*user.Filter, error) {
	switch userSession.Role() {
	case user.RoleAdmin:
		return requestedFilter, nil

	case user.RoleNetManager, user.RoleSelfOwner:
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

	case user.RoleManager:
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

	case user.RoleStaff, user.RoleUser:
		return nil, domainErr.NewForbiddenError("insufficient permissions to list users").
			WithDetail("role", userSession.Role().String())

	default:
		return nil, domainErr.NewForbiddenError("unknown role").
			WithDetail("role", userSession.Role().String())
	}
}
