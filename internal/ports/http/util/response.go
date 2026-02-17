package util

import (
	"context"
	"net/http"
)

func SendError(ctx context.Context, w http.ResponseWriter, err error) {}
