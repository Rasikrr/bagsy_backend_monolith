package schedule

import (
	"time"

	"github.com/google/uuid"
)

type EmployeeScheduleSlot struct {
	ID         uuid.UUID
	EmployeeID uuid.UUID
	Date       time.Time
	Type       SlotType
	StartTime  time.Time
	EndTime    time.Time
	CreatedAt  time.Time
	UpdatedAt  *time.Time
}

func NewEmployeeScheduleSlot(
	employeeID uuid.UUID,
	date time.Time,
	slotType SlotType,
	startTime, endTime time.Time,
) (*EmployeeScheduleSlot, error) {
	if !slotType.IsValid() {
		return nil, ErrInvalidSlotType
	}

	if !endTime.After(startTime) {
		return nil, ErrInvalidTimeRange
	}

	return &EmployeeScheduleSlot{
		ID:         uuid.New(),
		EmployeeID: employeeID,
		Date:       truncateToDate(date),
		Type:       slotType,
		StartTime:  startTime,
		EndTime:    endTime,
		CreatedAt:  time.Now(),
	}, nil
}

func (s *EmployeeScheduleSlot) IsWorkSlot() bool {
	return s.Type == SlotTypeWork
}

func (s *EmployeeScheduleSlot) IsRestSlot() bool {
	return s.Type == SlotTypeRest
}

func (s *EmployeeScheduleSlot) Duration() time.Duration {
	return s.EndTime.Sub(s.StartTime)
}

func (s *EmployeeScheduleSlot) Overlaps(other *EmployeeScheduleSlot) bool {
	if !sameDay(s.Date, other.Date) {
		return false
	}
	return s.StartTime.Before(other.EndTime) && other.StartTime.Before(s.EndTime)
}

func (s *EmployeeScheduleSlot) touch() {
	now := time.Now()
	s.UpdatedAt = &now
}
