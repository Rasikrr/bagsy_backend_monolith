package identity

import "errors"

var (
	ErrPhoneRequired         = errors.New("phone is required")
	ErrEmployeeNotFound      = errors.New("employee not found")
	ErrEmployeeDeleted       = errors.New("employee is deleted")
	ErrEmployeeInactive      = errors.New("employee is inactive")
	ErrEmployeeNameRequired  = errors.New("employee name is required")
	ErrEmployeePhoneRequired = errors.New("employee phone is required")
	ErrPermissionDenied      = errors.New("permission denied")
	ErrInvalidRole           = errors.New("invalid employee role")
	ErrInvalidGender         = errors.New("invalid gender")
)

var (
	ErrCustomerNotFound      = errors.New("customer not found")
	ErrCustomerDeleted       = errors.New("customer is deleted")
	ErrCustomerPhoneRequired = errors.New("customer phone is required")
	ErrCustomerNameRequired  = errors.New("customer name is required")
)

var (
	ErrCustomerBaseNotFound = errors.New("customer base record not found")
)
