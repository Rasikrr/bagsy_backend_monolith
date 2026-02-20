package catalog

import "errors"

var (
	ErrServiceNotFound         = errors.New("service not found")
	ErrEmployeeServiceNotFound = errors.New("employee service not found")
	ErrServiceDeleted          = errors.New("service deleted")
	ErrServiceNameRequired     = errors.New("service name required")
)

var (
	ErrCategorySelfParent   = errors.New("location_category self parent")
	ErrCategoryNameRequired = errors.New("location_category name is required")
)
