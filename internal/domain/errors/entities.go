package errors

// Bagsies errors
var (
	ErrBagsyNotFound              = NewNotFoundError("bagsies(ies) not found", nil)
	ErrBagsyTimeIsAlreadyOccupied = NewConflictError("time is already occupied", nil)
)

// Users errors
var (
	ErrUserNotFound      = NewNotFoundError("user(s) not found", nil)
	ErrUserActivated     = NewConflictError("user(s) activated", nil)
	ErrUserAlreadyExists = NewConflictError("user(s) already exists", nil)
)

// Networks errors
var (
	ErrNetworkNotFound      = NewNotFoundError("network(s) not found", nil)
	ErrNetworkAlreadyExists = NewConflictError("network(s) already exists", nil)
)

// Point Categories errors
var (
	ErrPointCategoryNotFound = NewNotFoundError("point category(ies) not found", nil)
	ErrPointAlreadyExists    = NewInvalidInputError("point with same code already exists", nil)
)

// Service categories errors
var (
	ErrServiceCategoryNotFound = NewNotFoundError("service category(ies) not found", nil)
)

// Service sub-categories errors
var (
	ErrServiceSubcategoryNotFound = NewNotFoundError("service subcategory(ies) not found", nil)
)

// Points errors
var (
	ErrPointNotFound = NewNotFoundError("point(s) not found", nil)
)

// Services errors
var (
	ErrServiceNotFound = NewNotFoundError("service(s) not found", nil)
)

// Master Services errors
var (
	ErrMasterServiceNotFound = NewNotFoundError("master service(s) not found", nil)
)

// Point Category Services errors
var (
	ErrPointCategoryServiceNotFound = NewNotFoundError("point category service(s) not found", nil)
)
