package tokens

import (
	"context"
	"errors"
	"fmt"
	"time"

	authDomain "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/google/uuid"
)

const (
	refreshTokenPrefix = "refresh_token:"
)

type RefreshTokenRepository struct {
	client *redis.Client
}

func NewRefreshTokenRepository(client *redis.Client) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		client: client,
	}
}

func (r *RefreshTokenRepository) SaveToken(ctx context.Context, tokenHash string, userID uuid.UUID, ttl time.Duration) error {
	key := r.makeKey(tokenHash)
	if err := r.client.SetWithExpiration(ctx, key, userID.String(), ttl); err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) GetToken(ctx context.Context, tokenHash string) (string, error) {
	key := r.makeKey(tokenHash)

	data, err := r.client.GetString(ctx, key)
	if errors.Is(err, redis.Nil) {
		return "", authDomain.ErrRefreshTokenNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get refresh token: %w", err)
	}

	return data, nil
}

func (r *RefreshTokenRepository) DeleteToken(ctx context.Context, tokenHash string) error {
	key := r.makeKey(tokenHash)
	if err := r.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) makeKey(tokenHash string) string {
	return refreshTokenPrefix + tokenHash
}
