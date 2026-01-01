package users

import (
	"context"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

type usersRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	ExistsByPhone(ctx context.Context, phone string) (bool, error)
}

type Service struct {
	usersRepo usersRepository
}

func NewService(
	usersRepo usersRepository,
) *Service {
	return &Service{
		usersRepo: usersRepo,
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

func (s *Service) Update(ctx context.Context, user *entity.User) error {
	return nil
}
