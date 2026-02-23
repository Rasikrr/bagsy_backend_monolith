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
	userSessionsPrefix = "user_sessions:"
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

	sessionKey := r.makeSessionKey(userID)
	if err := r.client.SAdd(ctx, sessionKey, tokenHash); err != nil {
		return fmt.Errorf("add token to user sessions: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) GetToken(ctx context.Context, tokenHash string) (uuid.UUID, error) {
	key := r.makeKey(tokenHash)

	data, err := r.client.GetBytes(ctx, key)
	if errors.Is(err, redis.Nil) {
		return uuid.Nil, authDomain.ErrRefreshTokenNotFound
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("get refresh token: %w", err)
	}

	return uuid.ParseBytes(data)
}

func (r *RefreshTokenRepository) DeleteToken(ctx context.Context, tokenHash string) error {
	// Get userID before deleting, so we can clean up the session set.
	userID, err := r.GetToken(ctx, tokenHash)
	if err != nil && !errors.Is(err, authDomain.ErrRefreshTokenNotFound) {
		return fmt.Errorf("get token for cleanup: %w", err)
	}

	key := r.makeKey(tokenHash)
	if err = r.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	if userID != uuid.Nil {
		sessionKey := r.makeSessionKey(userID)
		_ = r.client.SRem(ctx, sessionKey, tokenHash)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteAllByUserID(ctx context.Context, userID uuid.UUID) error {
	sessionKey := r.makeSessionKey(userID)

	hashes, err := r.client.SMembers(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("get user sessions: %w", err)
	}

	for _, hash := range hashes {
		key := r.makeKey(hash)
		_ = r.client.Delete(ctx, key)
	}

	if err = r.client.Delete(ctx, sessionKey); err != nil {
		return fmt.Errorf("delete user sessions set: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) makeKey(tokenHash string) string {
	return refreshTokenPrefix + tokenHash
}

func (r *RefreshTokenRepository) makeSessionKey(userID uuid.UUID) string {
	return userSessionsPrefix + userID.String()
}
