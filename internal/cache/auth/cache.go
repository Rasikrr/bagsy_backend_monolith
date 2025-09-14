package auth

import (
	"context"
	"fmt"
	"github.com/Rasikrr/bugsy_backend_monolith/internal/util/codegen"
	"github.com/Rasikrr/core/redis"
)

const (
	AuthCodePrefix = "codegen"
)

type Cache interface {
	Save(ctx context.Context, phone string)
}

type cache struct {
	cli redis.Cache
}

func NewCache(cli redis.Cache) Cache {
	return &cache{
		cli: cli,
	}
}

func (c *cache) Save(ctx context.Context, phone string) {
	code := codegen.GenerateAuthCode()

	err := c.cli.Set(ctx, fmt.Sprintf("%s:%s", AuthCodePrefix, phone), code)
	if err != nil {
		return
	}
}
