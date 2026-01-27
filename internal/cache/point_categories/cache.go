package pointcategories

import (
	"context"
	"encoding/json"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/point"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/cockroachdb/errors"
)

const cacheKey = "point_categories:all"

type Cache struct {
	cli *redis.Client
}

func New(cli *redis.Client) *Cache {
	return &Cache{cli: cli}
}

func (c *Cache) GetCategories(ctx context.Context) ([]*point.Category, error) {
	data, err := c.cli.GetString(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, domainErr.NewInternalError("failed to get point categories from cache", err)
	}

	var dtos []categoryDTO
	if err = json.Unmarshal([]byte(data), &dtos); err != nil {
		return nil, domainErr.NewInternalError("failed to unmarshal point categories from cache", err)
	}

	categories := make([]*point.Category, len(dtos))
	for i, dto := range dtos {
		categories[i] = dto.toDomain()
	}

	return categories, nil
}

func (c *Cache) SetCategories(ctx context.Context, categories []*point.Category, ttl time.Duration) error {
	dtos := make([]categoryDTO, len(categories))
	for i, cat := range categories {
		dtos[i] = toCategoryDTO(cat)
	}

	data, err := json.Marshal(dtos)
	if err != nil {
		return domainErr.NewInternalError("failed to marshal point categories for cache", err)
	}

	if err = c.cli.SetWithExpiration(ctx, cacheKey, string(data), ttl); err != nil {
		return domainErr.NewInternalError("failed to save point categories to cache", err)
	}

	return nil
}

type categoryDTO struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *string    `json:"updated_by,omitempty"`
}

func toCategoryDTO(c *point.Category) categoryDTO {
	return categoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		UpdatedBy:   c.UpdatedBy,
	}
}

func (d categoryDTO) toDomain() *point.Category {
	return &point.Category{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		UpdatedBy:   d.UpdatedBy,
	}
}
