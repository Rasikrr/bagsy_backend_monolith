package tokens

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	return r.client.SetWithExpiration(ctx, key, userID.String(), ttl)
}

func (r *RefreshTokenRepository) GetToken(ctx context.Context, tokenHash string) (string, error) {
	key := r.makeKey(tokenHash)

	data, err := r.client.GetString(ctx, key)
	if errors.Is(err, redis.Nil) {
		return "", fmt.Errorf("refresh token not found")
	}
	if err != nil {
		return "", fmt.Errorf("get refresh token: %w", err)
	}

	return data, nil
}

func (r *RefreshTokenRepository) DeleteToken(ctx context.Context, tokenHash string) error {
	key := r.makeKey(tokenHash)
	return r.client.Delete(ctx, key)
}

func (r *RefreshTokenRepository) makeKey(tokenHash string) string {
	return refreshTokenPrefix + tokenHash
}
