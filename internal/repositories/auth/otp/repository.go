package otp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/google/uuid"
)

const (
	otpPrefix = "appointment_otp:"
)

type Repository struct {
	client *redis.Client
}

func NewRepository(client *redis.Client) *Repository {
	return &Repository{
		client: client,
	}
}

func (r *Repository) Save(ctx context.Context, appointmentID uuid.UUID, o *auth.OTPCode) error {
	key := r.makeKey(appointmentID)

	m := fromDomain(appointmentID, o)
	data, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("marshal otp: %w", err)
	}

	ttl := time.Until(o.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("otp already expired")
	}

	if err := r.client.SetWithExpiration(ctx, key, data, ttl); err != nil {
		return fmt.Errorf("save otp to cache: %w", err)
	}
	return nil
}

func (r *Repository) GetByAppointmentID(ctx context.Context, appointmentID uuid.UUID) (*auth.OTPCode, error) {
	key := r.makeKey(appointmentID)

	data, err := r.client.GetBytes(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, fmt.Errorf("otp not found")
		}
		return nil, fmt.Errorf("get otp from cache: %w", err)
	}

	var m model
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal otp: %w", err)
	}

	return m.toDomain()
}

func (r *Repository) Delete(ctx context.Context, appointmentID uuid.UUID) error {
	key := r.makeKey(appointmentID)
	if err := r.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete otp from cache: %w", err)
	}
	return nil
}

func (r *Repository) makeKey(appointmentID uuid.UUID) string {
	return otpPrefix + appointmentID.String()
}
