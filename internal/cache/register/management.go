package register

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

func (c *Cache) SaveManagementRequest(ctx context.Context, req *command.RegisterManagementCommand) error {
	key := genManagementKey(req.Phone)
	return saveToCache(ctx, c, key, toRegisterManagementDTO(req))
}

func (c *Cache) DeleteManagementRequest(ctx context.Context, phone string) error {
	key := genManagementKey(phone)
	if err := c.cli.Delete(ctx, key); err != nil {
		return domainErr.NewInternalError("failed to delete register request", err)
	}
	return nil
}

func (c *Cache) GetManagementRequest(ctx context.Context, phone string) (*command.RegisterManagementCommand, error) {
	dto, err := getFromCache[registerManagementDTO](ctx, c, genManagementKey(phone))
	if err != nil {
		return nil, err
	}
	return dto.toDomain(), nil
}

func genManagementKey(phone string) string {
	return fmt.Sprintf("auth:management:confirm:%s", phone)
}
