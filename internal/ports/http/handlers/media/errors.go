package media

import (
	"net/http"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	httputil "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/util"
)

var mediaErrors = httputil.ErrorMap{
	media.ErrInvalidFileSize:     {Code: http.StatusBadRequest, Message: "file_size_limit"},
	media.ErrEmptyFilename:       {Code: http.StatusBadRequest, Message: "empty_filename"},
	media.ErrEmptyMimeType:       {Code: http.StatusBadRequest, Message: "empty_mime_type"},
	media.ErrUnsupportedMimeType: {Code: http.StatusBadRequest, Message: "unsupported_mime_type"},
	media.ErrInvalidPurpose:      {Code: http.StatusBadRequest, Message: "invalid_purpose"},
	media.ErrAssetNotFound:       {Code: http.StatusNotFound, Message: "asset_not_found"},
	media.ErrAssetNotPending:     {Code: http.StatusConflict, Message: "asset_not_pending"},
	media.ErrS3ObjectNotFound:    {Code: http.StatusNotFound, Message: "file_not_uploaded"},
}
