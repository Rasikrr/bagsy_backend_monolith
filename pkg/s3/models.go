package s3

import "time"

// ObjectInfo содержит информацию об объекте в S3
type ObjectInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
	ETag         string
}

// UploadPolicyOptions содержит параметры для генерации POST-формы загрузки
type UploadPolicyOptions struct {
	Key              string
	ContentType      string
	ContentLengthMin int64
	ContentLengthMax int64
	Expires          time.Duration
}

// UploadPolicyResponse содержит URL и поля формы для загрузки файла фронтендом
type UploadPolicyResponse struct {
	URL    string
	Fields map[string]string
}
