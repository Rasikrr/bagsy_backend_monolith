package errors

import (
	"context"
	"encoding/json"
	"net/http"

	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/Rasikrr/core/log"
	"github.com/cockroachdb/errors"
)

func HandleError(ctx context.Context, w http.ResponseWriter, err error) {
	logError(ctx, err)

	resp := toHTTPResponse(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Code)
	bb, _ := json.Marshal(resp)
	w.Write(bb)

}

func logError(ctx context.Context, err error) {
	var domErr *domainErr.DomainError
	if errors.As(err, &domErr) {
		status := getHTTPStatus(domErr.Type)
		if status >= http.StatusInternalServerError {
			log.Error(ctx, "internal server error occurred", log.Err(err))
		} else {
			log.Warn(ctx, "client error occurred", log.Err(err))
		}
	} else {
		log.Error(ctx, "unexpected error occurred", log.Err(err))
	}
}
