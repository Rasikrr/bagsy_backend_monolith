package errors

import "github.com/cockroachdb/errors"

// Bagsies errors
var (
	ErrBagsyNotFound = errors.New("bagsy(ies) not found")
)

// Users errors
var (
	ErrUserNotFound = errors.New("user(s) not found")
)

// Networks errors
var (
	ErrNetworkNotFound = errors.New("network(s) not found")
)

// Point Categories errors
var (
	ErrPointCategoryNotFound = errors.New("point category(ies) not found")
)

// Service categories errors
var (
	ErrServiceCategoryNotFound = errors.New("service category(ies) not found")
)

// Service sub-categories errors
var (
	ErrServiceSubcategoryNotFound = errors.New("service subcategory(ies) not found")
)

// Points errors
var (
	ErrPointNotFound = errors.New("point(s) not found")
)

// Services errors
var (
	ErrServiceNotFound = errors.New("service(s) not found")
)

// Master Services errors
var (
	ErrMasterServiceNotFound = errors.New("master service(s) not found")
)
