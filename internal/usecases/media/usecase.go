package media

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/media"
	"github.com/Rasikrr/bagsy_backend_monolith/pkg/s3"
)

type mediaRepository interface {
	Save(ctx context.Context, asset *media.Asset) error
}

type storageClient interface {
	GeneratePresignedPostURL(ctx context.Context, options s3.UploadPolicyOptions) (*s3.UploadPolicyResponse, error)
	BucketName() string
}

type UseCase struct {
	mediaRepo        mediaRepository
	storage          storageClient
	uploadExpires    time.Duration
	maxFileSizeBytes int64
}

func NewUseCase(
	mediaRepo mediaRepository,
	storage storageClient,
	uploadExpires time.Duration,
	maxFileSizeBytes int64,
) *UseCase {
	return &UseCase{
		mediaRepo:        mediaRepo,
		storage:          storage,
		uploadExpires:    uploadExpires,
		maxFileSizeBytes: maxFileSizeBytes,
	}
}
