package sms

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Rasikrr/core/redis"
)

const (
	spamPrefix = "SMS_SPAM"
)

type Cache interface {
	IsSpam(ctx context.Context, phone, message string) (bool, error)
	Set(ctx context.Context, phone, message string) error
}

type cache struct {
	spamTTL time.Duration
	cli     redis.Cache
}

func NewCache(cli redis.Cache, ttl time.Duration) Cache {
	return &cache{
		spamTTL: ttl,
		cli:     cli,
	}
}

func (c *cache) IsSpam(ctx context.Context, phone, message string) (bool, error) {
	out, err := c.cli.GetBool(ctx, c.genKey(phone, message))
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return false, nil
		}
		return false, err
	}
	return out, nil
}

func (c *cache) Set(ctx context.Context, phone, message string) error {
	return c.cli.SetWithExpiration(ctx, c.genKey(phone, message), true, c.spamTTL)
}

func (c *cache) genKey(phone, message string) string {
	return fmt.Sprintf("%s:%s:%s", spamPrefix, phone, message)
}
