package bagsyconfirm

import (
	"context"
	"fmt"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type Cache struct {
	cli     *redis.Client
	codeTTL time.Duration
}

func NewCache(cli *redis.Client, codeTTL time.Duration) *Cache {
	return &Cache{
		cli:     cli,
		codeTTL: codeTTL,
	}
}

func (c *Cache) GetCode(ctx context.Context, id uuid.UUID) (string, error) {
	authCode, err := c.cli.GetString(ctx, genKey(id))
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", domainErr.NewNotFoundError("code is not found or expired", err)
		}
		return "", domainErr.NewInternalError("error getting code", err)
	}
	return authCode, nil
}

func (c *Cache) SetCode(ctx context.Context, id uuid.UUID, code string) error {
	err := c.cli.SetWithExpiration(ctx, genKey(id), code, c.codeTTL)
	if err != nil {
		return domainErr.NewInternalError("error setting code", err)
	}
	return nil
}

func genKey(bagsyID uuid.UUID) string {
	return fmt.Sprintf("bagsy_confirm:%s", bagsyID.String())
}
