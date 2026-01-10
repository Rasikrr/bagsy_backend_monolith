package media

import (
	"context"
	"time"
)

type storage interface {
	GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

type Service struct {
	storage storage
}

// nolint
func NewService(storage storage) *Service {
	return &Service{
		storage: storage,
	}
}
