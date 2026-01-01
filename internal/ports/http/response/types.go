package response

type EmptySuccessResponse struct {
	Message string `json:"message"`
}

func NewEmptySuccessResponse(message string) *EmptySuccessResponse {
	return &EmptySuccessResponse{
		Message: message,
	}
}
