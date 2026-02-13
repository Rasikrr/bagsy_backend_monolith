package notification

import "errors"

var (
	ErrTaskLocked   = errors.New("notification task is locked by another worker")
	ErrTaskNotFound = errors.New("notification task not found")
)
