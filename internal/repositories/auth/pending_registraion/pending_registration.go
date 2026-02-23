package pendingregistraion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
	"github.com/Rasikrr/core/cache/redis"
)

const (
	pendingRegPrefix = "pending_reg:"
)

type PendingRegistrationStore struct {
	client *redis.Client
}

func NewPendingRegistrationStore(client *redis.Client) *PendingRegistrationStore {
	return &PendingRegistrationStore{
		client: client,
	}
}

func (s *PendingRegistrationStore) Save(ctx context.Context, reg *auth.PendingRegistration) error {
	key := s.makeKey(reg.Phone)

	m := fromPendingRegistration(reg)
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal pending registration: %w", err)
	}

	ttl := time.Until(reg.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("pending registration already expired")
	}

	if err = s.client.SetWithExpiration(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("save pending registration: %w", err)
	}
	return nil
}

func (s *PendingRegistrationStore) Get(ctx context.Context, phone shared.Phone) (*auth.PendingRegistration, error) {
	key := s.makeKey(phone)

	data, err := s.client.GetBytes(ctx, key)
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get pending registration: %w", err)
	}

	var m pendingRegistrationModel
	if err = json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal pending registration: %w", err)
	}

	return m.toUseCase()
}

func (s *PendingRegistrationStore) Delete(ctx context.Context, phone shared.Phone) error {
	key := s.makeKey(phone)
	if err := s.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete pending registration: %w", err)
	}
	return nil
}

func (s *PendingRegistrationStore) makeKey(phone shared.Phone) string {
	return pendingRegPrefix + phone.String()
}
