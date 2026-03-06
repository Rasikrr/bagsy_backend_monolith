package catalog

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type ServiceCategoryTree struct {
	ID        uuid.UUID
	Name      string
	SortOrder int
	Children  []ServiceCategoryTree
}

func (u *UseCase) GetServiceCategories(ctx context.Context, locationCategoryID uuid.UUID) ([]ServiceCategoryTree, error) {
	categories, err := u.catalogRepo.GetServiceCategoriesByLocationCategoryID(ctx, locationCategoryID)
	if err != nil {
		return nil, errors.Wrap(err, "get service categories")
	}

	childrenByParent := make(map[uuid.UUID][]ServiceCategoryTree)
	var roots []ServiceCategoryTree

	for _, c := range categories {
		node := ServiceCategoryTree{
			ID:        c.ID,
			Name:      c.Name,
			SortOrder: c.SortOrder,
		}
		if c.ParentID == nil {
			roots = append(roots, node)
		} else {
			childrenByParent[*c.ParentID] = append(childrenByParent[*c.ParentID], node)
		}
	}

	for i := range roots {
		if children, ok := childrenByParent[roots[i].ID]; ok {
			roots[i].Children = children
		}
	}

	return roots, nil
}
