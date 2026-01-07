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

// ConvertAlmatyTimeToUTC parses time string in Almaty timezone and returns time.Time in Almaty timezone
// Example: "09:00" -> 0001-01-01 09:00:00 +0500 Asia/Almaty
func ConvertAlmatyTimeToUTC(timeStr string) (time.Time, error) {
	// Almaty is UTC+5 (fixed offset)
	almaty := time.FixedZone("Asia/Almaty", 5*60*60)

	// Use zero date (0001-01-01) to indicate time-only value
	baseDate := "0001-01-01"
	fullTime := baseDate + " " + timeStr

	t, err := time.ParseInLocation("2006-01-02 15:04", fullTime, almaty)
	if err != nil {
		return time.Time{}, err
	}

	// Return time in Almaty timezone (not UTC)
	return t, nil
}
