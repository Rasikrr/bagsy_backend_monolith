package pending

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	"github.com/Rasikrr/core/cache/redis"
)

const (
	pendingInvitePrefix = "pending_invite:"
)

type Store struct {
	client *redis.Client
}

func NewStore(client *redis.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Save(ctx context.Context, inv *invite.PendingInvite) error {
	key := s.makeKey(inv.Phone)

	m := fromPendingInvite(inv)
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal pending invite: %w", err)
	}

	ttl := time.Until(inv.ExpiresAt)
	if ttl <= 0 {
		return authDomain.ErrInviteTokenExpired
	}

	if err = s.client.SetWithExpiration(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("save pending invite: %w", err)
	}
	return nil
}

func (s *Store) Get(ctx context.Context, phone shared.Phone) (*invite.PendingInvite, error) {
	key := s.makeKey(phone)

	data, err := s.client.GetBytes(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil, authDomain.ErrInviteTokenNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get pending invite: %w", err)
	}

	var m pendingInviteModel
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal pending invite: %w", err)
	}

	return m.toUseCase()
}

func (s *Store) Delete(ctx context.Context, phone shared.Phone) error {
	key := s.makeKey(phone)
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete pending invite: %w", err)
	}
	return nil
}

func (s *Store) makeKey(phone shared.Phone) string {
	return pendingInvitePrefix + phone.String()
}
