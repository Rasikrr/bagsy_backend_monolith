package s3

import "github.com/cockroachdb/errors"

// Validation errors
var (
	ErrEmptyRegion    = errors.New("s3: AWS region is required")
	ErrEmptyAccessKey = errors.New("s3: AWS access key is required")
	ErrEmptySecretKey = errors.New("s3: AWS secret key is required")
	ErrEmptyBucket    = errors.New("s3: bucket name is required")
	ErrEmptyKey       = errors.New("s3: object key is required")
	ErrEmptyData      = errors.New("s3: data to upload cannot be empty")
	ErrNilReader      = errors.New("s3: reader cannot be nil")
	ErrNilWriter      = errors.New("s3: writer cannot be nil")
	ErrEmptyKeys      = errors.New("s3: keys cannot be empty")
	ErrNoValidKeys    = errors.New("s3: no valid keys provided")
	ErrInvalidExpiry  = errors.New("s3: expiration time must be positive")
)

// Operation errors
var (
	ErrConfigFailed   = errors.New("s3: failed to load AWS configuration")
	ErrUploadFailed   = errors.New("s3: failed to upload file")
	ErrDownloadFailed = errors.New("s3: failed to download file")
	ErrDeleteFailed   = errors.New("s3: failed to delete file")
	ErrListFailed     = errors.New("s3: failed to list objects")
	ErrEmptyLocation  = errors.New("s3: empty location returned after upload")
	ErrCheckFailed    = errors.New("s3: failed to check object existence")
	ErrPresignFailed  = errors.New("s3: failed to generate presigned URL")
)
