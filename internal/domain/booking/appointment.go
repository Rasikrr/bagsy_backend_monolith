package booking

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

// reasonAutoCompleted — причина в истории статусов при системном авто-завершении.
const reasonAutoCompleted = "auto-completed by system"

// ─────────────────────────────────────────────────────────────────
// Aggregate Root: Appointment
// ─────────────────────────────────────────────────────────────────

type Appointment struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	LocationID     uuid.UUID
	ServiceID      uuid.UUID
	EmployeeID     uuid.UUID
	CustomerID     uuid.UUID

	StartAt time.Time
	EndAt   time.Time

	Price           shared.Money
	DurationMinutes shared.Duration

	Status             Status
	CustomerComment    *string
	CancelledBy        *uuid.UUID
	CancellationReason *string

	StatusHistory []StatusHistoryEntry

	CreatedAt time.Time
	UpdatedAt *time.Time
}

type CreateAppointmentParams struct {
	OrganizationID  uuid.UUID
	LocationID      uuid.UUID
	ServiceID       uuid.UUID
	EmployeeID      uuid.UUID
	CustomerID      uuid.UUID
	StartAt         time.Time
	DurationMinutes shared.Duration
	Price           shared.Money
	CustomerComment *string
}

func NewAppointment(params CreateAppointmentParams) (*Appointment, error) {
	endAt := params.StartAt.Add(params.DurationMinutes.AsDuration())
	now := time.Now()

	appointment := &Appointment{
		ID:              uuid.New(),
		OrganizationID:  params.OrganizationID,
		LocationID:      params.LocationID,
		ServiceID:       params.ServiceID,
		EmployeeID:      params.EmployeeID,
		CustomerID:      params.CustomerID,
		StartAt:         params.StartAt,
		EndAt:           endAt,
		Price:           params.Price,
		DurationMinutes: params.DurationMinutes,
		Status:          StatusPending,
		CustomerComment: params.CustomerComment,
		StatusHistory:   make([]StatusHistoryEntry, 0),
		CreatedAt:       now,
	}

	appointment.addStatusHistory(StatusPending, nil, nil)

	return appointment, nil
}

// ─────────────────────────────────────────────────────────────────
// Business Methods
// ─────────────────────────────────────────────────────────────────

func (a *Appointment) Confirm(by uuid.UUID) error {
	if !a.Status.CanTransitionTo(StatusConfirmed) {
		return ErrAppointmentInvalidStatusTransition
	}

	a.Status = StatusConfirmed
	a.addStatusHistory(StatusConfirmed, &by, nil)
	a.touch()

	return nil
}

func (a *Appointment) Start(by uuid.UUID) error {
	if !a.Status.CanTransitionTo(StatusInProgress) {
		return ErrAppointmentInvalidStatusTransition
	}

	a.Status = StatusInProgress
	a.addStatusHistory(StatusInProgress, &by, nil)
	a.touch()

	return nil
}

func (a *Appointment) Complete(by uuid.UUID) error {
	if !a.Status.CanTransitionTo(StatusCompleted) {
		return ErrAppointmentInvalidStatusTransition
	}

	a.Status = StatusCompleted
	a.addStatusHistory(StatusCompleted, &by, nil)
	a.touch()

	return nil
}

// AutoComplete завершает прошедшую запись системно (по расписанию, без
// действующего сотрудника). Допустимо из confirmed/in_progress: время визита
// прошло, поэтому запись считается выполненной независимо от того, перевёл ли
// её кто-то вручную в in_progress.
func (a *Appointment) AutoComplete() error {
	if a.Status != StatusConfirmed && a.Status != StatusInProgress {
		return ErrAppointmentInvalidStatusTransition
	}

	a.Status = StatusCompleted
	a.addStatusHistory(StatusCompleted, nil, ptr(reasonAutoCompleted)) // changedBy=nil — системное действие
	a.touch()

	return nil
}

func (a *Appointment) Cancel(by uuid.UUID, reason string) error {
	if !a.Status.CanTransitionTo(StatusCancelled) {
		return ErrAppointmentInvalidStatusTransition
	}

	a.Status = StatusCancelled
	a.CancelledBy = &by
	a.CancellationReason = &reason
	a.addStatusHistory(StatusCancelled, &by, &reason)
	a.touch()

	return nil
}

func (a *Appointment) Reschedule(newStart time.Time, by uuid.UUID) error {
	if a.IsFinal() {
		return ErrAppointmentIsFinal
	}
	newEnd := newStart.Add(a.DurationMinutes.AsDuration())

	if !newEnd.After(newStart) {
		return ErrInvalidTimeRange
	}

	if newStart.Before(time.Now()) {
		return ErrCannotScheduleInPast
	}

	a.StartAt = newStart
	a.EndAt = newEnd
	a.addStatusHistory(StatusConfirmed, &by, ptr("rescheduled"))
	a.touch()

	return nil
}

func (a *Appointment) UpdateComment(comment *string) error {
	if a.IsFinal() {
		return ErrAppointmentIsFinal
	}

	a.CustomerComment = comment
	a.touch()

	return nil
}

// ─────────────────────────────────────────────────────────────────
// Query Methods
// ─────────────────────────────────────────────────────────────────

func (a *Appointment) IsPending() bool {
	return a.Status == StatusPending
}

func (a *Appointment) IsConfirmed() bool {
	return a.Status == StatusConfirmed
}

func (a *Appointment) IsInProgress() bool {
	return a.Status == StatusInProgress
}

func (a *Appointment) IsCompleted() bool {
	return a.Status == StatusCompleted
}

func (a *Appointment) IsCancelled() bool {
	return a.Status == StatusCancelled
}

func (a *Appointment) IsFinal() bool {
	return a.Status.IsFinal()
}

func (a *Appointment) IsUpcoming() bool {
	return !a.IsFinal() && a.StartAt.After(time.Now())
}

func (a *Appointment) IsOngoing() bool {
	now := time.Now()
	return !a.IsFinal() && a.StartAt.Before(now) && a.EndAt.After(now)
}

func (a *Appointment) IsPast() bool {
	return a.EndAt.Before(time.Now())
}

func (a *Appointment) BelongsTo(organizationID uuid.UUID) bool {
	return a.OrganizationID == organizationID
}

func (a *Appointment) BelongsToLocation(locationID uuid.UUID) bool {
	return a.LocationID == locationID
}

func (a *Appointment) BelongsToEmployee(employeeID uuid.UUID) bool {
	return a.EmployeeID == employeeID
}

func (a *Appointment) Duration() time.Duration {
	return a.DurationMinutes.AsDuration()
}

func (a *Appointment) Overlaps(start, end time.Time) bool {
	return a.StartAt.Before(end) && start.Before(a.EndAt)
}

// ─────────────────────────────────────────────────────────────────
// Private Methods
// ─────────────────────────────────────────────────────────────────

func (a *Appointment) touch() {
	now := time.Now()
	a.UpdatedAt = &now
}

func (a *Appointment) addStatusHistory(toStatus Status, changedBy *uuid.UUID, reason *string) {
	var fromStatus *Status
	if len(a.StatusHistory) > 0 {
		last := a.StatusHistory[len(a.StatusHistory)-1]
		fromStatus = &last.ToStatus
	}

	a.StatusHistory = append(a.StatusHistory, StatusHistoryEntry{
		ID:         uuid.New(),
		FromStatus: fromStatus,
		ToStatus:   toStatus,
		ChangedBy:  changedBy,
		Reason:     reason,
		CreatedAt:  time.Now(),
	})
}

func ptr(s string) *string {
	return &s
}
