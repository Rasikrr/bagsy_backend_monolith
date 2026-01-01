package errors

type ErrorResponse struct {
	Message string         `json:"message"`
	Code    int            `json:"code"`
	Details map[string]any `json:"details,omitempty"`
}
