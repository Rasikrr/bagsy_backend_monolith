package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/Rasikrr/core/redis"
)

const (
	codeKey = "code"
)

func (c *cache) SetCode(ctx context.Context, phone, code string) error {
	err := c.cli.SetWithExpiration(ctx, c.genKey(c.genCodeKey(phone)), code, c.codeTTL)
	if err != nil {
		return fmt.Errorf("failed to set code to cache: %w", err)
	}
	return nil
}

func (c *cache) GetCode(ctx context.Context, phone string) (string, error) {
	code, err := c.cli.Get(ctx, c.genKey(c.genCodeKey(phone)))
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			return "", errors.New("code not found")
		}
		return "", errors.New("failed to get code from cache")
	}
	codeStr, ok := code.(string)
	if !ok {
		return "", errors.New("invalid code type")
	}
	return codeStr, nil
}

func (c *cache) CheckSpam(ctx context.Context, phone string) (bool, error) {
	res, err := c.cli.Exists(ctx, c.genKey(c.genCodeKey(phone)))
	if err != nil {
		return false, err
	}
	return res, nil
}

func (c *cache) genCodeKey(phone string) string {
	return fmt.Sprintf("%s:%s", codeKey, phone)
}
