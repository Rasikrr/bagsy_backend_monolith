package media

import "strings"

var allowedMimeTypes = map[string]struct{}{
	"image/jpeg":      {},
	"image/png":       {},
	"image/webp":      {},
	"video/mp4":       {},
	"video/quicktime": {},
	"video/webm":      {},
}

// MimeType — Value Object для допустимых MIME-типов медиафайлов.
type MimeType struct {
	value string
}

func ParseMimeType(raw string) (MimeType, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return MimeType{}, ErrEmptyMimeType
	}
	if _, ok := allowedMimeTypes[raw]; !ok {
		return MimeType{}, ErrUnsupportedMimeType
	}
	return MimeType{value: raw}, nil
}

func (m MimeType) String() string { return m.value }

func (m MimeType) IsImage() bool {
	return strings.HasPrefix(m.value, "image/")
}

func (m MimeType) IsVideo() bool {
	return strings.HasPrefix(m.value, "video/")
}
