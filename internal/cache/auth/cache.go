package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/Rasikrr/core/redis"
)

const (
	authPrefix = "AUTH"
)

type Cache interface {
	SetCode(ctx context.Context, phone, code string) error
	CheckSpam(ctx context.Context, phone string) (bool, error)
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

func (c *cache) genKey(key string) string {
	return fmt.Sprintf("%s:%s", authPrefix, key)
}
