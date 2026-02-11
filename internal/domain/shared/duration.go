package shared

import "time"

type Duration struct {
	minutes int
}

func NewDuration(val int) (Duration, error) {
	if val < 0 {
		return Duration{}, ErrInvalidDuration
	}
	return Duration{minutes: val}, nil
}

func (d Duration) Minutes() int {
	return d.minutes
}

func (d Duration) AsDuration() time.Duration {
	return time.Duration(d.minutes) * time.Minute
}
