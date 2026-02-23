package booking

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/booking"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type appointmentModel struct {
	ID                 uuid.UUID       `db:"id"`
	OrganizationID     uuid.UUID       `db:"organization_id"`
	LocationID         uuid.UUID       `db:"location_id"`
	ServiceID          uuid.UUID       `db:"service_id"`
	EmployeeID         uuid.UUID       `db:"employee_id"`
	CustomerID         uuid.UUID       `db:"customer_id"`
	StartAt            time.Time       `db:"start_at"`
	EndAt              time.Time       `db:"end_at"`
	Price              decimal.Decimal `db:"price"`
	DurationMinutes    int             `db:"duration_minutes"`
	Status             string          `db:"status"`
	CustomerComment    *string         `db:"customer_comment"`
	CancelledBy        *uuid.UUID      `db:"cancelled_by"`
	CancellationReason *string         `db:"cancellation_reason"`
	CreatedAt          time.Time       `db:"created_at"`
	UpdatedAt          *time.Time      `db:"updated_at"`
}

type statusHistoryModel struct {
	ID            uuid.UUID  `db:"id"`
	AppointmentID uuid.UUID  `db:"appointment_id"`
	FromStatus    *string    `db:"from_status"`
	ToStatus      string     `db:"to_status"`
	Payload       []byte     `db:"payload"`
	ChangedBy     *uuid.UUID `db:"changed_by"`
	Reason        *string    `db:"reason"`
	CreatedAt     time.Time  `db:"created_at"`
}

func fromDomain(a *booking.Appointment) *appointmentModel {
	return &appointmentModel{
		ID:                 a.ID,
		OrganizationID:     a.OrganizationID,
		LocationID:         a.LocationID,
		ServiceID:          a.ServiceID,
		EmployeeID:         a.EmployeeID,
		CustomerID:         a.CustomerID,
		StartAt:            a.StartAt,
		EndAt:              a.EndAt,
		Price:              a.Price.Amount(),
		DurationMinutes:    a.DurationMinutes.Minutes(),
		Status:             string(a.Status),
		CustomerComment:    a.CustomerComment,
		CancelledBy:        a.CancelledBy,
		CancellationReason: a.CancellationReason,
		CreatedAt:          a.CreatedAt,
		UpdatedAt:          a.UpdatedAt,
	}
}

func (m *appointmentModel) toDomain(history []booking.StatusHistoryEntry) (*booking.Appointment, error) {
	price, err := shared.NewMoney(m.Price)
	if err != nil {
		return nil, err
	}

	duration, err := shared.NewDuration(m.DurationMinutes)
	if err != nil {
		return nil, err
	}

	return &booking.Appointment{
		ID:                 m.ID,
		OrganizationID:     m.OrganizationID,
		LocationID:         m.LocationID,
		ServiceID:          m.ServiceID,
		EmployeeID:         m.EmployeeID,
		CustomerID:         m.CustomerID,
		StartAt:            m.StartAt,
		EndAt:              m.EndAt,
		Price:              price,
		DurationMinutes:    duration,
		Status:             booking.Status(m.Status),
		CustomerComment:    m.CustomerComment,
		CancelledBy:        m.CancelledBy,
		CancellationReason: m.CancellationReason,
		StatusHistory:      history,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
	}, nil
}

func fromHistoryDomain(appointmentID uuid.UUID, h booking.StatusHistoryEntry) *statusHistoryModel {
	var fromStatus *string
	if h.FromStatus != nil {
		s := string(*h.FromStatus)
		fromStatus = &s
	}

	return &statusHistoryModel{
		ID:            h.ID,
		AppointmentID: appointmentID,
		FromStatus:    fromStatus,
		ToStatus:      string(h.ToStatus),
		ChangedBy:     h.ChangedBy,
		Reason:        h.Reason,
		CreatedAt:     h.CreatedAt,
	}
}

type calendarEntryRow struct {
	AppointmentID   uuid.UUID       `db:"appointment_id"`
	Status          string          `db:"status"`
	StartAt         time.Time       `db:"start_at"`
	EndAt           time.Time       `db:"end_at"`
	Price           decimal.Decimal `db:"price"`
	DurationMinutes int             `db:"duration_minutes"`
	CustomerComment *string         `db:"customer_comment"`

	EmployeeID   uuid.UUID `db:"employee_id"`
	EmployeeName string    `db:"employee_name"`

	CustomerID    uuid.UUID `db:"customer_id"`
	CustomerName  string    `db:"customer_name"`
	CustomerPhone string    `db:"customer_phone"`

	ServiceID    uuid.UUID `db:"service_id"`
	ServiceName  string    `db:"service_name"`
	ServiceColor string    `db:"service_color"`

	LocationID   uuid.UUID `db:"location_id"`
	LocationName string    `db:"location_name"`
}

func (m *statusHistoryModel) toDomain() booking.StatusHistoryEntry {
	var fromStatus *booking.Status
	if m.FromStatus != nil {
		s := booking.Status(*m.FromStatus)
		fromStatus = &s
	}

	return booking.StatusHistoryEntry{
		ID:         m.ID,
		FromStatus: fromStatus,
		ToStatus:   booking.Status(m.ToStatus),
		ChangedBy:  m.ChangedBy,
		Reason:     m.Reason,
		CreatedAt:  m.CreatedAt,
	}
}
