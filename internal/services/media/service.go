package media

import (
	"context"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/clients/s3"
	domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
	"github.com/cockroachdb/errors"
)

type storage interface {
	GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, expiresIn time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, expiresIn time.Duration) (string, error)
}

type Service struct {
	storage  storage
	mediaTTL time.Duration
}

func NewService(
	storage storage,
	mediaTTL time.Duration,
) *Service {
	return &Service{
		storage:  storage,
		mediaTTL: mediaTTL,
	}
}

func (s *Service) GenerateUploadURL(ctx context.Context, key string, contentType string) (string, error) {
	url, err := s.storage.GeneratePresignedUploadURL(ctx, key, contentType, s.mediaTTL)
	if err != nil {
		return "", s.mapS3Error(err)
	}
	return url, nil
}

func (s *Service) GenerateDownloadURL(ctx context.Context, key string) (string, error) {
	url, err := s.storage.GeneratePresignedDownloadURL(ctx, key, s.mediaTTL)
	if err != nil {
		return "", s.mapS3Error(err)
	}
	return url, nil
}

// mapS3Error преобразует ошибки S3 клиента → доменные ошибки
func (s *Service) mapS3Error(err error) *domainErr.DomainError {
	if err == nil {
		return nil
	}

	// ========== VALIDATION ERRORS (400) ==========

	// Required fields validation
	if errors.Is(err, s3.ErrEmptyRegion) {
		return domainErr.NewInvalidInputError("AWS region is required", err)
	}
	if errors.Is(err, s3.ErrEmptyAccessKey) {
		return domainErr.NewInvalidInputError("AWS access key is required", err)
	}
	if errors.Is(err, s3.ErrEmptySecretKey) {
		return domainErr.NewInvalidInputError("AWS secret key is required", err)
	}
	if errors.Is(err, s3.ErrEmptyBucket) {
		return domainErr.NewInvalidInputError("bucket name is required", err)
	}
	if errors.Is(err, s3.ErrEmptyKey) {
		return domainErr.NewInvalidInputError("object key is required", err)
	}
	if errors.Is(err, s3.ErrEmptyData) {
		return domainErr.NewInvalidInputError("data to upload cannot be empty", err)
	}
	if errors.Is(err, s3.ErrNilReader) {
		return domainErr.NewInvalidInputError("reader cannot be nil", err)
	}
	if errors.Is(err, s3.ErrNilWriter) {
		return domainErr.NewInvalidInputError("writer cannot be nil", err)
	}
	if errors.Is(err, s3.ErrEmptyKeys) {
		return domainErr.NewInvalidInputError("keys cannot be empty", err)
	}
	if errors.Is(err, s3.ErrNoValidKeys) {
		return domainErr.NewInvalidInputError("no valid keys provided", err)
	}
	if errors.Is(err, s3.ErrInvalidExpiry) {
		return domainErr.NewInvalidInputError("expiration time must be positive", err)
	}

	// ========== OPERATION ERRORS (500) ==========

	// Configuration errors
	if errors.Is(err, s3.ErrConfigFailed) {
		return domainErr.NewInternalError("failed to configure storage service", err)
	}

	// Upload/Download errors
	if errors.Is(err, s3.ErrUploadFailed) {
		return domainErr.NewInternalError("failed to upload file to storage", err)
	}
	if errors.Is(err, s3.ErrDownloadFailed) {
		return domainErr.NewInternalError("failed to download file from storage", err)
	}
	if errors.Is(err, s3.ErrDeleteFailed) {
		return domainErr.NewInternalError("failed to delete file from storage", err)
	}
	if errors.Is(err, s3.ErrListFailed) {
		return domainErr.NewInternalError("failed to list files in storage", err)
	}
	if errors.Is(err, s3.ErrCheckFailed) {
		return domainErr.NewInternalError("failed to check file existence in storage", err)
	}

	// Presign errors
	if errors.Is(err, s3.ErrPresignFailed) {
		return domainErr.NewInternalError("failed to generate file access URL", err)
	}

	// Empty location
	if errors.Is(err, s3.ErrEmptyLocation) {
		return domainErr.NewInternalError("storage returned empty file location", err)
	}

	// ========== FALLBACK ==========

	// Для всех неизвестных ошибок
	return domainErr.NewInternalError("storage service error", err)
}
