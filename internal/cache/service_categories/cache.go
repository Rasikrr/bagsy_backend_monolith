package servicecategories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/service"
	"github.com/Rasikrr/core/cache/redis"
	"github.com/cockroachdb/errors"
)

const cacheKeyPrefix = "service_categories:point_category:"

type Cache struct {
	cli *redis.Client
}

func New(cli *redis.Client) *Cache {
	return &Cache{cli: cli}
}

func (c *Cache) Get(ctx context.Context, pointCategoryID int) ([]*service.CategoryWithSubcategories, error) {
	key := fmt.Sprintf("%s%d", cacheKeyPrefix, pointCategoryID)
	data, err := c.cli.GetString(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}
		return nil, domainErr.NewInternalError("failed to get service categories from cache", err)
	}

	var dtos []categoryDTO
	if err = json.Unmarshal([]byte(data), &dtos); err != nil {
		return nil, domainErr.NewInternalError("failed to unmarshal service categories from cache", err)
	}

	result := make([]*service.CategoryWithSubcategories, len(dtos))
	for i, dto := range dtos {
		result[i] = dto.toDomain()
	}

	return result, nil
}

func (c *Cache) Set(ctx context.Context, pointCategoryID int, categories []*service.CategoryWithSubcategories, ttl time.Duration) error {
	key := fmt.Sprintf("%s%d", cacheKeyPrefix, pointCategoryID)

	dtos := make([]categoryDTO, len(categories))
	for i, cat := range categories {
		dtos[i] = toCategoryDTO(cat)
	}

	data, err := json.Marshal(dtos)
	if err != nil {
		return domainErr.NewInternalError("failed to marshal service categories for cache", err)
	}

	if err = c.cli.SetWithExpiration(ctx, key, string(data), ttl); err != nil {
		return domainErr.NewInternalError("failed to save service categories to cache", err)
	}

	return nil
}

type categoryDTO struct {
	ID            int              `json:"id"`
	Name          string           `json:"name"`
	Description   *string          `json:"description,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     *time.Time       `json:"updated_at,omitempty"`
	UpdatedBy     *string          `json:"updated_by,omitempty"`
	Subcategories []subcategoryDTO `json:"subcategories"`
}

type subcategoryDTO struct {
	ID          int        `json:"id"`
	CategoryID  int        `json:"category_id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
	UpdatedBy   *string    `json:"updated_by,omitempty"`
}

func toCategoryDTO(c *service.CategoryWithSubcategories) categoryDTO {
	subDTOs := make([]subcategoryDTO, len(c.Subcategories))
	for i, sub := range c.Subcategories {
		subDTOs[i] = subcategoryDTO{
			ID:          sub.ID,
			CategoryID:  sub.CategoryID,
			Name:        sub.Name,
			Description: sub.Description,
			CreatedAt:   sub.CreatedAt,
			UpdatedAt:   sub.UpdatedAt,
			UpdatedBy:   sub.UpdatedBy,
		}
	}
	return categoryDTO{
		ID:            c.Category.ID,
		Name:          c.Category.Name,
		Description:   c.Category.Description,
		CreatedAt:     c.Category.CreatedAt,
		UpdatedAt:     c.Category.UpdatedAt,
		UpdatedBy:     c.Category.UpdatedBy,
		Subcategories: subDTOs,
	}
}

func (d categoryDTO) toDomain() *service.CategoryWithSubcategories {
	subs := make([]*service.Subcategory, len(d.Subcategories))
	for i, subDTO := range d.Subcategories {
		subs[i] = &service.Subcategory{
			ID:          subDTO.ID,
			CategoryID:  subDTO.CategoryID,
			Name:        subDTO.Name,
			Description: subDTO.Description,
			CreatedAt:   subDTO.CreatedAt,
			UpdatedAt:   subDTO.UpdatedAt,
			UpdatedBy:   subDTO.UpdatedBy,
		}
	}
	return &service.CategoryWithSubcategories{
		Category: &service.Category{
			ID:          d.ID,
			Name:        d.Name,
			Description: d.Description,
			CreatedAt:   d.CreatedAt,
			UpdatedAt:   d.UpdatedAt,
			UpdatedBy:   d.UpdatedBy,
		},
		Subcategories: subs,
	}
}
