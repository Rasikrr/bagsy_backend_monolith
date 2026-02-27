package media

import (
	"net/http"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	coreHTTP "github.com/Rasikrr/core/http"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// confirm handles POST /api/v1/media/{id}/confirm.
//
// @Summary      Подтверждение загрузки файла
// @Description  Проверяет наличие файла в S3 и меняет статус ассета на uploaded.
// @Tags         media
// @Produce      json
// @Param        id   path  string  true  "Asset ID"
// @Success      204
// @Failure      404  {object}  httputil.errorResponse  "asset_not_found / file_not_uploaded"
// @Failure      409  {object}  httputil.errorResponse  "asset_not_pending"
// @Failure      500  {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/media/{id}/confirm [post]
func (h *Handler) confirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	if err = h.mediaUseCase.ConfirmUpload(ctx, id); err != nil {
		httputil.SendError(ctx, w, err, mediaErrors)
		return
	}

	coreHTTP.SendData(ctx, w, nil, http.StatusNoContent)
}
