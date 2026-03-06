package notification

import "errors"

var (
	ErrTaskLocked     = errors.New("notification task is locked by another worker")
	ErrTaskNotFound   = errors.New("notification task not found")
	ErrTaskNotPending = errors.New("notification task is not in pending status")
)
