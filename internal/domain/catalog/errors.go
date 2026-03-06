package catalog

import "errors"

var (
	ErrServiceNotFound         = errors.New("service not found")
	ErrEmployeeServiceNotFound = errors.New("employee service not found")
	ErrServiceDeleted          = errors.New("service deleted")
	ErrServiceInactive         = errors.New("service is inactive")
	ErrServiceNameRequired     = errors.New("service name required")
	ErrServiceInvalidColor     = errors.New("service invalid color")
)

var (
	ErrCategorySelfParent       = errors.New("location_category self parent")
	ErrCategoryNameRequired     = errors.New("location_category name is required")
	ErrServiceCategoryNotFound  = errors.New("service category not found")
	ErrCategoryMismatch         = errors.New("service category does not match location category")
	ErrEmployeeLocationMismatch = errors.New("employee does not belong to service location")
)
