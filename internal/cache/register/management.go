package register

import (
	"context"
	"fmt"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	authS "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"
)

func (c *Cache) SaveManagementRequest(ctx context.Context, req *authS.ManagementRegistrationState, ttl time.Duration) error {
	key := genManagementKey(req.Command.Phone)
	return saveToCache(ctx, c, key, toRegisterManagementDTO(req), ttl)
}

func (c *Cache) DeleteManagementRequest(ctx context.Context, phone string) error {
	key := genManagementKey(phone)
	if err := c.cli.Delete(ctx, key); err != nil {
		return domainErr.NewInternalError("failed to delete register request", err)
	}
	return nil
}

func (c *Cache) GetManagementRequest(ctx context.Context, phone string) (*authS.ManagementRegistrationState, error) {
	dto, err := getFromCache[registerManagementStateDTO](ctx, c, genManagementKey(phone))
	if err != nil {
		return nil, err
	}
	return dto.toDomain(), nil
}

func genManagementKey(phone string) string {
	return fmt.Sprintf("auth:management:confirm:%s", phone)
}
