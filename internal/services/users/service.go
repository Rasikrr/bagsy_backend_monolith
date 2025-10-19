package users

import (
	"context"
	"errors"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	appErrors "github.com/Rasikrr/bagsy_backend_monolith/internal/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
	"github.com/samber/lo"
)

type Service interface {
	Create(ctx context.Context, user *entity.User) error

	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	GetByPointCode(ctx context.Context, pointCode string) ([]*entity.User, error)
	GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.User, error)

	ExistsByPhone(ctx context.Context, phone string) (bool, error)

	Update(ctx context.Context, phone string, params UpdateParams) error
	SetPasswordByPhone(ctx context.Context, phone string, password string) error
	SetActive(ctx context.Context, phone string) error

	DeleteUnactivatedUsers(ctx context.Context, olderThan time.Duration) error
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
	existingUser, err := s.usersRepo.GetByPhone(ctx, user.Phone)
	if err != nil && !errors.Is(err, appErrors.ErrUserNotFound) {
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

func (s *service) GetByPointCode(ctx context.Context, pointCode string) ([]*entity.User, error) {
	out, err := s.usersRepo.GetByParams(ctx, users.GetParams{
		PointCode: &pointCode,
	})
	if err != nil {
		return nil, errGetUser.Wrap(err)
	}
	return out, nil
}

func (s *service) GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.User, error) {
	out, err := s.usersRepo.GetByParams(ctx, users.GetParams{
		NetworkCode: &networkCode,
	})
	if err != nil {
		return nil, errGetUser.Wrap(err)
	}
	return out, nil
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
