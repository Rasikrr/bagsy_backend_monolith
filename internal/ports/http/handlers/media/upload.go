package media

import (
	"net/http"

	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
	uc "github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/media"
	coreHTTP "github.com/Rasikrr/core/http"
)

// upload handles POST /api/v1/media/upload.
//
// @Summary      Генерация presigned URL для загрузки файла
// @Description  Создаёт запись о медиафайле и возвращает presigned POST URL для загрузки в S3.
// @Description  Допустимые значения purpose: avatars, organizations, locations, services, service-categories.
// @Tags         media
// @Accept       json
// @Produce      json
// @Param        body  body      uploadRequest  true  "Параметры файла"
// @Success      201   {object}  uploadResponse
// @Failure      400   {object}  httputil.errorResponse  "invalid_purpose / file_size_limit / empty_filename / unsupported_mime_type"
// @Failure      500   {object}  httputil.errorResponse
// @Security     ApiKeyAuth
// @Router       /api/v1/media/upload [post]
func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req uploadRequest
	if err := coreHTTP.GetData(r, &req); err != nil {
		httputil.SendBadRequest(ctx, w, err)
		return
	}

	input := uc.GenerateUploadURLInput{
		Filename:  req.Filename,
		MimeType:  req.MimeType,
		SizeBytes: req.SizeBytes,
		Purpose:   req.Purpose,
	}

	out, err := h.mediaUseCase.GenerateUploadURL(ctx, input)
	if err != nil {
		httputil.SendError(ctx, w, err, mediaErrors)
		return
	}

	coreHTTP.SendData(ctx, w, uploadResponse{
		AssetID:      out.AssetID.String(),
		UploadURL:    out.UploadURL,
		UploadFields: out.UploadFields,
	}, http.StatusCreated)
}
