package media

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/request"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/response"
)

// getUploadURL godoc
// @Summary Получить пресайнд ссылку для загрузки медиафайла
// @Description Создает запись в БД со статусом PENDING и генерирует временную S3 URL для загрузки файла напрямую с клиента
// @Tags media
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body getUploadURLRequest true "Данные о загружаемом файле"
// @Success 200 {object} getUploadURLResponse "Ссылки на загрузку и метаданные"
// @Failure 400 {object} errors.ErrorResponse "Неверные параметры запроса или неподдерживаемый тип контента"
// @Failure 401 {object} errors.ErrorResponse "Пользователь не авторизован"
// @Failure 500 {object} errors.ErrorResponse "Внутренняя ошибка сервера"
// @Router /api/v1/media/upload [post]
func (c *Controller) getUploadURL(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req getUploadURLRequest
	if err := request.GetAndValidateData(r, &req); err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	out, err := c.mediaService.GenerateUploadURL(ctx, req.Filename, req.ContentType, req.Purpose)
	if err != nil {
		errors.HandleError(ctx, w, err)
		return
	}
	response.SendData(ctx, w, toUploadURLResponse(out), http.StatusOK)
}
