package errors

import (
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
)

// HTTPStatusMap maps domain error types to HTTP status codes
var httpStatusMap = map[domainErr.ErrorType]int{
	domainErr.TypeNotFound:        http.StatusNotFound,            // 404
	domainErr.TypeInvalidInput:    http.StatusBadRequest,          // 400
	domainErr.TypeValidation:      http.StatusBadRequest,          // 400
	domainErr.TypeUnauthorized:    http.StatusUnauthorized,        // 401
	domainErr.TypeForbidden:       http.StatusForbidden,           // 403
	domainErr.TypeConflict:        http.StatusConflict,            // 409
	domainErr.TypeTooManyRequests: http.StatusTooManyRequests,     // 429
	domainErr.TypeInternal:        http.StatusInternalServerError, // 500
}

// GetHTTPStatus returns HTTP status code for domain error type
// Returns 500 (Internal Server Error) if error type is unknown
func getHTTPStatus(errType domainErr.ErrorType) int {
	if status, ok := httpStatusMap[errType]; ok {
		return status
	}
	return http.StatusInternalServerError
}

func toHTTPResponse(err error) ErrorResponse {
	var (
		domErr *domainErr.DomainError
	)
	if errors.As(err, &domErr) {
		httpCode := getHTTPStatus(domErr.Type)

		return ErrorResponse{
			Message: domErr.Message,
			Code:    httpCode,
			Details: domErr.Details,
		}
	}
	return ErrorResponse{
		Message: "internal server error",
		Code:    http.StatusInternalServerError,
	}
}
