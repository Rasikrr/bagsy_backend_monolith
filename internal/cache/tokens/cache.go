package tokens

import (
	"context"
	"fmt"
	"time"

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

// DeleteToken удаляет токен (logout)
func (c *Cache) DeleteRefreshToken(ctx context.Context, tokenHash string) error {
	key := c.refreshTokenKey(tokenHash)
	if err := c.cli.Delete(ctx, key); err != nil {
		return domainErr.NewInternalError("failed to delete refresh token", err)
	}
	return nil
}

// MarkRegistrationTokenAsUsed помечает registration token как использованный (one-time use)
// Использует SET NX для атомарности - предотвращает race condition при одновременных запросах
// Возвращает true если токен уже был использован ранее
func (c *Cache) MarkRegistrationTokenAsUsed(ctx context.Context, tokenHash string, ttl time.Duration) (alreadyUsed bool, err error) {
	key := c.registrationTokenKey(tokenHash)

	// SetNX - атомарная операция "set if not exists"
	// ok = true  → ключ установлен (первое использование токена)
	// ok = false → ключ уже существует (повторное использование)
	ok, err := c.cli.SetNX(ctx, key, "used", ttl)
	if err != nil {
		return false, domainErr.NewInternalError("failed to mark registration token as used", err)
	}

	// Инвертируем результат для понятности:
	// ok=true  → первый раз → alreadyUsed=false
	// ok=false → уже был использован → alreadyUsed=true
	return !ok, nil
}

func (c *Cache) refreshTokenKey(tokenHash string) string {
	return fmt.Sprintf("refresh_token:%s", tokenHash)
}

func (c *Cache) registrationTokenKey(tokenHash string) string {
	return fmt.Sprintf("registration_token:%s", tokenHash)
}
