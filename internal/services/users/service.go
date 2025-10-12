package users

import (
	"context"
	"errors"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	appErrors "github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/session"
	"github.com/samber/lo"
)

type Service interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByParams(ctx context.Context, params GetParams) ([]*entity.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
	SetPasswordByPhone(ctx context.Context, phone string, password string) error
	SetActive(ctx context.Context, phone string) error
	DeleteUnactivatedUsers(ctx context.Context, olderThan time.Duration) error
	Update(ctx context.Context, phone string, params UpdateParams) error
}

type service struct {
	usersRepo users.Repository
}

func NewService(usersRepo users.Repository) Service {
	return &service{
		usersRepo: usersRepo,
	}
}

func (s *service) Create(ctx context.Context, user *entity.User) error {
	by, err := session.GetSession(ctx)
	if err != nil {
		return appErrors.ErrSessionNotFound
	}
	if !by.Role.HasPermission(user.Role) {
		return errNoPermission
	}
	networkCode := by.GetNetworkCode()
	user.NetworkCode = &networkCode

	existingUser, err := s.usersRepo.GetByPhone(ctx, user.Phone)
	if err != nil && errors.Is(err, appErrors.ErrUserNotFound) {
		return errCreateUser.Wrap(err)
	}
	if existingUser.Active {
		return errUserAlreadyExists
	}
	if createErr := s.usersRepo.Create(ctx, user); createErr != nil {
		return errCreateUser.Wrap(createErr)
	}
	return nil
}

func (s *service) GetByParams(ctx context.Context, params GetParams) ([]*entity.User, error) {
	by, err := session.GetSession(ctx)
	if err != nil {
		return nil, appErrors.ErrSessionNotFound
	}
	err = params.validate(by)
	if err != nil {
		return nil, errValidateParams.Wrap(err)
	}
	user, err := s.usersRepo.GetByParams(ctx, params.convert())
	if err != nil {
		return nil, errGetUser.Wrap(err)
	}
	return user, nil
}

func (s *service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, errGetUser.Wrap(err)
	}
	return user, nil
}

func (s *service) SetPasswordByPhone(ctx context.Context, phone string, password string) error {
	patch := users.NewUserUpdatePatch().
		SetPhones(phone).
		SetPassword(password).
		Build()
	err := s.usersRepo.Update(ctx, patch)
	if err != nil {
		return errSetPassword.Wrap(err)
	}
	return nil
}

func (s *service) ExistsByPhone(ctx context.Context, phone string) (bool, error) {
	return s.usersRepo.ExistsByPhone(ctx, phone)
}

func (s *service) SetActive(ctx context.Context, phone string) error {
	patch := users.NewUserUpdatePatch().
		SetPhones(phone).
		SetActive(true).
		Build()
	err := s.usersRepo.Update(ctx, patch)
	if err != nil {
		return errActivateUser.Wrap(err)
	}
	return nil
}

func (s *service) Update(ctx context.Context, phone string, params UpdateParams) error {
	patch := params.ToPatch(phone)
	err := s.usersRepo.Update(ctx, patch)
	if err != nil {
		return errUpdateUser.Wrap(err)
	}
	return nil
}

func (s *service) DeleteUnactivatedUsers(ctx context.Context, olderThan time.Duration) error {
	userToDelete, err := s.usersRepo.GetInactive(ctx, olderThan)
	if err != nil {
		return errDeleteUnactivatedUsers.Wrap(err)
	}
	if len(userToDelete) == 0 {
		return nil
	}
	phones := lo.Map(userToDelete, func(u *entity.User, _ int) string {
		return u.Phone
	})
	err = s.usersRepo.SoftDelete(ctx, phones...)
	if err != nil {
		return errDeleteUnactivatedUsers.Wrap(err)
	}
	return nil
}
