package schedule

import "errors"

var (
	ErrInvalidSlotType  = errors.New("invalid slot type")
	ErrInvalidTimeRange = errors.New("invalid time range")
	ErrSlotNotFound     = errors.New("schedule slot not found")
)
