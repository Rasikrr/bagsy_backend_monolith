package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/core/redis"
)

const (
	prefix     = "auth"
	codePrefix = "code"
)

type Cache interface {
	Save(ctx context.Context, phone string) error
}

type cache struct {
	codeTTL time.Duration
	cli     redis.Cache
}

func NewCache(cli redis.Cache, codeTTL time.Duration) Cache {
	return &cache{
		cli:     cli,
		codeTTL: codeTTL,
	}
}

func (c *cache) Save(ctx context.Context, phone string) error {
	code := codegen.GenerateAuthCode()

	err := c.cli.SetWithExpiration(ctx, c.codeKey(phone), code, c.codeTTL)
	if err != nil {
		return fmt.Errorf("failed to set code to cache: %v", err)
	}
	return nil
}

func (c *cache) codeKey(phone string) string {
	return fmt.Sprintf("%s:%s:%s", prefix, codePrefix, phone)
}
