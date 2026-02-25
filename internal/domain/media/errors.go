package media

import "errors"

var (
	ErrInvalidFileSize     = errors.New("file size limit")
	ErrEmptyFilename       = errors.New("filename could not be empty")
	ErrEmptyMimeType       = errors.New("mime-type could not be empty")
	ErrUnsupportedMimeType = errors.New("unsupported mime-type")
	ErrInvalidPurpose      = errors.New("invalid media purpose")
)
