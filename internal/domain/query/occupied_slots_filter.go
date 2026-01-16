package query

import "time"

type OccupiedSlotsFilter struct {
	PointCode    string
	MasterPhones []string
	StartAt      time.Time
	EndAt        time.Time
}
