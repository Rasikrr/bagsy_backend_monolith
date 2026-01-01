package response

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	coreCtx "github.com/Rasikrr/core/context"
	coreHttp "github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
)

func SendData(ctx context.Context, w http.ResponseWriter, data interface{}, status int) {
	traceID, ok := coreCtx.TraceID(ctx)
	if ok {
		w.Header().Set(coreHttp.TraceIDHeader, traceID)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	marshaller, ok := data.(json.Marshaler)
	if ok {
		bb, err := marshaller.MarshalJSON()
		if err != nil {
			errors.HandleError(ctx, w, err)
			return
		}
		w.Write(bb)
		logSuccessResponse(ctx, len(bb), status)
		return
	}

	bb, err := json.Marshal(data)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	w.Write(bb)
	logSuccessResponse(ctx, len(bb), status)
}

func logSuccessResponse(ctx context.Context, bodyLen, status int) {
	log.Info(ctx,
		"successfully handled request",
		log.Int("response_length", bodyLen),
		log.Int("HTTP_status", status),
	)
}
