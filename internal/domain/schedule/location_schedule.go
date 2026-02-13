package schedule

import (
	"time"

	"github.com/google/uuid"
)

type LocationScheduleSlot struct {
	ID         uuid.UUID
	LocationID uuid.UUID
	Date       time.Time
	Type       SlotType
	StartTime  time.Time
	EndTime    time.Time
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewLocationScheduleSlot(
	locationID uuid.UUID,
	date time.Time,
	slotType SlotType,
	startTime, endTime time.Time,
) (*LocationScheduleSlot, error) {
	if !slotType.IsValid() {
		return nil, ErrInvalidSlotType
	}

	if !endTime.After(startTime) {
		return nil, ErrInvalidTimeRange
	}

	return &LocationScheduleSlot{
		ID:         uuid.New(),
		LocationID: locationID,
		Date:       truncateToDate(date),
		Type:       slotType,
		StartTime:  startTime,
		EndTime:    endTime,
		CreatedAt:  time.Now(),
	}, nil
}

func (s *LocationScheduleSlot) IsWorkSlot() bool {
	return s.Type == SlotTypeWork
}

func (s *LocationScheduleSlot) IsRestSlot() bool {
	return s.Type == SlotTypeRest
}

func (s *LocationScheduleSlot) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

func (s *LocationScheduleSlot) Overlaps(other *LocationScheduleSlot) bool {
	if !sameDay(s.Date, other.Date) {
		return false
	}
	return s.StartTime.Before(other.EndTime) && other.StartTime.Before(s.EndTime)
}
