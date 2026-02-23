package actiontoken

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/core/cache/redis"
)

const (
	actionTokenPrefix = "action_token:"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{client: client}
}

func (s *Store) Save(ctx context.Context, token *authDomain.ActionToken) error {
	key := s.makeKey(token.Token)
	m := fromDomain(token)

	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal action token: %w", err)
	}

	ttl := time.Until(token.ExpiresAt)
	if ttl <= 0 {
		return authDomain.ErrOTPExpired
	}

	if err = s.client.SetWithExpiration(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("save action token: %w", err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, token string) (*authDomain.ActionToken, error) {
	key := s.makeKey(token)

	data, err := s.client.GetBytes(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil, authDomain.ErrActionTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get action token: %w", err)
	}

	var m model
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal action token: %w", err)
	}

	return m.toDomain()
}

func (s *Store) Delete(ctx context.Context, token string) error {
	key := s.makeKey(token)
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete action token: %w", err)
	}
	return nil
}

func (s *Store) makeKey(token string) string {
	return actionTokenPrefix + token
}
