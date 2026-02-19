package token

import (
	"context"
	"errors"
	"fmt"
	"time"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/core/cache/redis"
)

const (
	inviteTokenPrefix = "invite_token:"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Save(ctx context.Context, token string, phone shared.Phone, ttl time.Duration) error {
	key := s.makeKey(token)
	if err := s.client.SetWithExpiration(ctx, key, phone.String(), ttl); err != nil {
		return fmt.Errorf("save invite token: %w", err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, token string) (shared.Phone, error) {
	key := s.makeKey(token)

	data, err := s.client.GetBytes(ctx, key)
	if errors.Is(err, redis.Nil) {
		return shared.Phone{}, authDomain.ErrInviteTokenNotFound
	}
	if err != nil {
		return shared.Phone{}, fmt.Errorf("get invite token: %w", err)
	}

	phone, err := shared.NewPhone(string(data))
	if err != nil {
		return shared.Phone{}, fmt.Errorf("parse phone from invite token: %w", err)
	}

	return phone, nil
}

func (s *Store) Delete(ctx context.Context, token string) error {
	key := s.makeKey(token)
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete invite token: %w", err)
	}
	return nil
}

func (s *Store) makeKey(token string) string {
	return inviteTokenPrefix + token
}
