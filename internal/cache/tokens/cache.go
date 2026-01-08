package tokens

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/dto"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/cockroachdb/errors"
)

type Cache struct {
	cli *redis.Client
}

func New(cli *redis.Client) *Cache {
	return &Cache{cli: cli}
}

// SaveRefreshToken сохраняет hash токена с метаданными
// Key: refresh_token:<hash>
// Value: phone
// TTL: время жизни токена
func (c *Cache) SaveRefreshToken(ctx context.Context, tokenHash, phone string, ttl time.Duration) error {
	key := c.refreshTokenKey(tokenHash)
	if err := c.cli.SetWithExpiration(ctx, key, phone, ttl); err != nil {
		return domainErr.NewInternalError("failed to save refresh token", err)
	}
	return nil
}

// GetRefreshToken получает phone по hash токена
func (c *Cache) GetRefreshToken(ctx context.Context, tokenHash string) (string, error) {
	key := c.refreshTokenKey(tokenHash)
	phone, err := c.cli.GetString(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domainErr.ErrInvalidToken.WithDetail("reason", "token not found or expired")
		}
		return "", domainErr.NewInternalError("failed to get refresh token", err)
	}
	return phone, nil
}

// DeleteRefreshToken удаляет токен (logout)
func (c *Cache) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	key := c.refreshTokenKey(tokenHash)
	if err := c.cli.Delete(ctx, key); err != nil {
		return domainErr.NewInternalError("failed to delete refresh token", err)
	}
	return nil
}

func (c *Cache) refreshTokenKey(tokenHash string) string {
	return fmt.Sprintf("refresh_token:%s", tokenHash)
}

// SaveAuthToken сохраняет короткий authтокен с данными регистрации
// Key: invite_token:<token>
// Value: JSON с phone, point_code, network_code
// TTL: время жизни токена (обычно 1 час)
func (c *Cache) SaveAuthToken(ctx context.Context, token string, payload *dto.AuthTokenPayload, ttl time.Duration) error {
	key := c.authTokenKey(token)

	// Сериализуем payload в JSON
	data, err := json.Marshal(payload)
	if err != nil {
		return domainErr.NewInternalError("failed to marshal auth token payload", err)
	}

	if err = c.cli.SetWithExpiration(ctx, key, string(data), ttl); err != nil {
		return domainErr.NewInternalError("failed to save auth token", err)
	}
	return nil
}

// GetAuthToken получает данные по короткому authтокену
func (c *Cache) GetAuthToken(ctx context.Context, token string) (*dto.AuthTokenPayload, error) {
	key := c.authTokenKey(token)
	data, err := c.cli.GetString(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domainErr.ErrInvalidToken.WithDetail("reason", "auth token not found or expired")
		}
		return nil, domainErr.NewInternalError("failed to get auth token", err)
	}

	// Десериализуем JSON
	var payload dto.AuthTokenPayload
	if err = json.Unmarshal([]byte(data), &payload); err != nil {
		return nil, domainErr.NewInternalError("failed to unmarshal auth token payload", err)
	}

	return &payload, nil
}

// DeleteAuthToken удаляет authтокен (например, после использования)
func (c *Cache) DeleteAuthToken(ctx context.Context, token string) error {
	key := c.authTokenKey(token)
	if err := c.cli.Delete(ctx, key); err != nil {
		return domainErr.NewInternalError("failed to delete auth token", err)
	}
	return nil
}

func (c *Cache) authTokenKey(token string) string {
	return fmt.Sprintf("auth_token:%s", token)
}
