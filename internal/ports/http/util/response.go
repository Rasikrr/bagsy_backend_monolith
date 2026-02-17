package util

import (
	"context"
	"net/http"

	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/Rasikrr/core/log"
	"github.com/cockroachdb/errors"
)

// ErrorInfo описывает HTTP-код и клиентское сообщение для доменной ошибки.
type ErrorInfo struct {
	Code    int
	Message string
}

// ErrorMap — маппинг доменных ошибок на HTTP-представление.
type ErrorMap map[error]ErrorInfo

// SendError ищет ошибку в переданной карте и отправляет соответствующий HTTP-ответ.
// Если ошибка не найдена в карте — логирует и отдаёт 500.
func SendError(ctx context.Context, w http.ResponseWriter, err error, em ErrorMap) {
	for target, info := range em {
		if errors.Is(err, target) {
			coreHTTP.SendData(ctx, w, errorResponse{Error: info.Message}, info.Code)
			return
		}
	}

	log.Error(ctx, "internal server error", log.Err(err))
	coreHTTP.SendData(ctx, w, errorResponse{Error: "internal_error"}, http.StatusInternalServerError)
}

// SendBadRequest отправляет 400 с переданным сообщением.
// Используется для ошибок парсинга/валидации запроса.
func SendBadRequest(ctx context.Context, w http.ResponseWriter, err error) {
	log.Warn(ctx, "bad request", log.Err(err))
	coreHTTP.SendData(ctx, w, errorResponse{Error: "bad_request"}, http.StatusBadRequest)
}

type errorResponse struct {
	Error string `json:"error"`
}
