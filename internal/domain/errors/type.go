package errors

import (
	"fmt"

	"github.com/cockroachdb/errors"
)

type ErrorType string

const (
	TypeNotFound        ErrorType = "NOT_FOUND"
	TypeInvalidInput    ErrorType = "INVALID_INPUT"
	TypeValidation      ErrorType = "VALIDATION"
	TypeUnauthorized    ErrorType = "UNAUTHORIZED"
	TypeForbidden       ErrorType = "FORBIDDEN"
	TypeConflict        ErrorType = "CONFLICT"
	TypeTooManyRequests ErrorType = "TOO_MANY_REQUESTS"
	TypeInternal        ErrorType = "INTERNAL"
)

type DomainError struct {
	Type    ErrorType
	Message string
	Cause   error
	Details map[string]interface{}
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *DomainError) Unwrap() error { return e.Cause }

func (e *DomainError) WithError(err error) *DomainError {
	e.Cause = err
	return e
}

func (e *DomainError) WithDetail(key string, value interface{}) *DomainError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// NewNotFoundError creates a new not found domain error
func NewNotFoundError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    TypeNotFound,
		Message: message,
		Cause:   cause,
	}
}

func NewInvalidInputError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    TypeInvalidInput,
		Message: message,
		Cause:   cause,
	}
}

func NewValidationError(message string, keyVals ...string) *DomainError {
	err := &DomainError{Type: TypeValidation, Message: message}
	if len(keyVals) > 0 {
		if len(keyVals)%2 != 0 {
			keyVals = append(keyVals, "no description")
		}
		err.Details = make(map[string]interface{})
		for i := 0; i < len(keyVals); i += 2 {
			err.Details[keyVals[i]] = keyVals[i+1]
		}
	}
	return err
}

func NewTooManyRequestsError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    TypeTooManyRequests,
		Message: message,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Type:    TypeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Type:    TypeForbidden,
		Message: message,
	}
}

func NewConflictError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    TypeConflict,
		Message: message,
		Cause:   cause,
	}
}

func NewInternalError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    TypeInternal,
		Message: message,
		Cause:   errors.WithStack(cause),
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == TypeNotFound
	}
	return false
}

func IsInvalidInput(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == TypeInvalidInput
	}
	return false
}

func IsValidation(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == TypeValidation
	}
	return false
}

func IsInternal(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == TypeInternal
	}
	return false
}

func IsConflict(err error) bool {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type == TypeConflict
	}
	return false
}

func GetType(err error) ErrorType {
	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		return domainErr.Type
	}
	return TypeInternal
}
