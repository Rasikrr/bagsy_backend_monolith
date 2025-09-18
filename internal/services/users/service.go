package users

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/domain/entity"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/repositories/users"
)

type Service interface {
	Create(ctx context.Context, user *entity.User) error
	GetByPhone(ctx context.Context, phone string) (*entity.User, error)
	SetPasswordByPhone(ctx context.Context, phone string, password string) error
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
	if err != nil && !errors.Is(err, users.ErrUserNotFound) {
		return fmt.Errorf("get user by phone: %w", err)
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}
	if createErr := s.usersRepo.Create(ctx, user); createErr != nil {
		return fmt.Errorf("create user: %w", createErr)
	}
	return nil
}

func (s *service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	user, err := s.usersRepo.GetByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("get user by phone: %w", err)
	}
	return user, nil
}

func (s *service) SetPasswordByPhone(ctx context.Context, phone string, password string) error {
	err := s.usersRepo.SetPassword(ctx, phone, password)
	if err != nil {
		return errors.New("cannot set password")
	}

	return nil
}
