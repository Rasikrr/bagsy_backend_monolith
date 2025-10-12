package time

import "time"

func ConvertStrToScheduleTime(s string) time.Time {
	// Assuming the input string is in "15:04" format
	parsedTime, err := time.Parse("15:04", s)
	if err != nil {
		return time.Time{}
	}
	return parsedTime
}
