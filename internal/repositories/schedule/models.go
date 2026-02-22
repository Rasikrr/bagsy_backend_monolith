package schedule

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/schedule"
	"github.com/google/uuid"
)

type locationSlotModel struct {
	ID         uuid.UUID  `db:"id"`
	LocationID uuid.UUID  `db:"location_id"`
	Date       time.Time  `db:"date"`
	Type       string     `db:"type"`
	StartTime  time.Time  `db:"start_time"`
	EndTime    time.Time  `db:"end_time"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (m *locationSlotModel) toDomain() (*schedule.LocationScheduleSlot, error) {
	st, err := schedule.ParseSlotType(m.Type)
	if err != nil {
		return nil, err
	}

	return &schedule.LocationScheduleSlot{
		ID:         m.ID,
		LocationID: m.LocationID,
		Date:       m.Date,
		Type:       st,
		StartTime:  m.StartTime,
		EndTime:    m.EndTime,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}

type employeeSlotModel struct {
	ID         uuid.UUID  `db:"id"`
	EmployeeID uuid.UUID  `db:"employee_id"`
	Date       time.Time  `db:"date"`
	Type       string     `db:"type"`
	StartTime  time.Time  `db:"start_time"`
	EndTime    time.Time  `db:"end_time"`
	CreatedAt  time.Time  `db:"created_at"`
	UpdatedAt  *time.Time `db:"updated_at"`
}

func (m *employeeSlotModel) toDomain() (*schedule.EmployeeScheduleSlot, error) {
	st, err := schedule.ParseSlotType(m.Type)
	if err != nil {
		return nil, err
	}

	return &schedule.EmployeeScheduleSlot{
		ID:         m.ID,
		EmployeeID: m.EmployeeID,
		Date:       m.Date,
		Type:       st,
		StartTime:  m.StartTime,
		EndTime:    m.EndTime,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}
