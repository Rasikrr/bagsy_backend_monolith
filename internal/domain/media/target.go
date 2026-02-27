package media

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Purpose string

const (
	PurposeAvatar          Purpose = "avatars"
	PurposeOrganization    Purpose = "organizations"
	PurposeLocation        Purpose = "locations"
	PurposeService         Purpose = "services"
	PurposeServiceCategory Purpose = "service-categories"
)

func (p Purpose) IsValid() bool {
	switch p {
	case PurposeAvatar, PurposeOrganization, PurposeLocation, PurposeService, PurposeServiceCategory:
		return true
	default:
		return false
	}
}

func ParsePurpose(raw string) (Purpose, error) {
	p := Purpose(raw)
	if !p.IsValid() {
		return p, ErrInvalidPurpose
	}
	return p, nil
}

// buildObjectKey генерирует flat S3 key: {purpose}/{YYYY/MM/DD}/{assetID}{ext}
func buildObjectKey(purpose Purpose, assetID uuid.UUID, filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	datePart := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s%s", purpose, datePart, assetID, ext)
}
