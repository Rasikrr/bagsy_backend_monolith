package media

import "errors"

var (
	ErrInvalidFileSize = errors.New("file size limit")
	ErrEmptyFilename   = errors.New("filename could not be empty")
	ErrEmptyMimeType   = errors.New("mime-type could not be empty")
)
