package register

import (
	"context"
	"encoding/json"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/cockroachdb/errors"
)

type Cache struct {
	cli *redis.Client
}

func NewCache(cli *redis.Client) *Cache {
	return &Cache{
		cli: cli,
	}
}

func saveToCache[T any](ctx context.Context, c *Cache, key string, dto T, ttl time.Duration) error {
	bb, err := json.Marshal(dto)
	if err != nil {
		return domainErr.NewInternalError("failed to marshal register request", err)
	}

	if err = c.cli.SetWithExpiration(ctx, key, bb, ttl); err != nil {
		return domainErr.NewInternalError("failed to save register request", err)
	}
	return nil
}

func getFromCache[T any](ctx context.Context, c *Cache, key string) (*T, error) {
	bb, err := c.cli.GetBytes(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, domainErr.NewNotFoundError("register timeout or not request found", err)
		}
		return nil, domainErr.NewInternalError("failed to fetch register request", err)
	}

	var dto T
	if err = json.Unmarshal(bb, &dto); err != nil {
		return nil, domainErr.NewInternalError("failed to unmarshal register request", err)
	}
	return &dto, nil
}
