package s3

import "time"

// ObjectInfo содержит информацию об объекте в S3
type ObjectInfo struct {
	Key          string
	Size         int64
	LastModified time.Time
	ETag         string
}
