package register

import (
	"context"
	"fmt"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/command"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

func (c *Cache) SaveStaffRequest(ctx context.Context, req *command.RegisterStaffCommand) error {
	key := genStaffKey(req.Phone)
	return saveToCache(ctx, c, key, toRegisterStaffDTO(req))
}

func (c *Cache) GetStaffRequest(ctx context.Context, phone string) (*command.RegisterStaffCommand, error) {
	dto, err := getFromCache[registerStaffDTO](ctx, c, genStaffKey(phone))
	if err != nil {
		return nil, err
	}
	return dto.toDomain(), nil
}

func (c *Cache) DeleteStaffRequest(ctx context.Context, phone string) error {
	key := genStaffKey(phone)
	if err := c.cli.Delete(ctx, key); err != nil {
		return domainErr.NewInternalError("failed to delete register request", err)
	}
	return nil
}

func genStaffKey(phone string) string {
	return fmt.Sprintf("auth:staff:confirm:%s", phone)
}
