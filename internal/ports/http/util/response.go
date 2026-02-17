package util

import (
	"context"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func SendError(ctx context.Context, w http.ResponseWriter, err error) {}
