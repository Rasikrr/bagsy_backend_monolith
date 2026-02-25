package media

import "github.com/google/uuid"

type GenerateUploadURLInput struct {
	Filename  string
	MimeType  string
	SizeBytes int64
	Purpose   string
}

type GenerateUploadURLOutput struct {
	AssetID      uuid.UUID
	UploadURL    string
	UploadFields map[string]string
}
