package catalog

import "errors"

var (
	ErrServiceDeleted      = errors.New("service deleted")
	ErrServiceNameRequired = errors.New("service name required")
)

var (
	ErrCategorySelfParent   = errors.New("category self parent")
	ErrCategoryNameRequired = errors.New("category name is required")
)
