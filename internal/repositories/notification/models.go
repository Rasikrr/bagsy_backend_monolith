package notification

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type taskModel struct {
	ID             int64      `db:"id"`
	AppointmentID  uuid.UUID  `db:"appointment_id"`
	Type           string     `db:"type"`
	RecipientType  string     `db:"recipient_type"`
	RecipientPhone string     `db:"recipient_phone"`
	Metadata       []byte     `db:"metadata"`
	Status         string     `db:"status"`
	ScheduledFor   time.Time  `db:"scheduled_for"`
	Attempts       int        `db:"attempts"`
	MaxAttempts    int        `db:"max_attempts"`
	LastError      *string    `db:"last_error"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      *time.Time `db:"updated_at"`
}

type metadataDTO struct {
	ServiceName   string    `json:"service_name"`
	LocationName  string    `json:"location_name"`
	AppointmentAt time.Time `json:"appointment_at"`
}

func (m *taskModel) toDomain() (*notification.Task, error) {
	phone, err := shared.NewPhone(m.RecipientPhone)
	if err != nil {
		return nil, fmt.Errorf("parse recipient phone: %w", err)
	}

	var dto metadataDTO
	if err = json.Unmarshal(m.Metadata, &dto); err != nil {
		return nil, fmt.Errorf("unmarshal notification metadata: %w", err)
	}

	return &notification.Task{
		ID:             m.ID,
		AppointmentID:  m.AppointmentID,
		Type:           notification.Type(m.Type),
		RecipientType:  notification.RecipientType(m.RecipientType),
		RecipientPhone: phone,
		Metadata: notification.Metadata{
			ServiceName:   dto.ServiceName,
			LocationName:  dto.LocationName,
			AppointmentAt: dto.AppointmentAt,
		},
		Status:       notification.Status(m.Status),
		ScheduledFor: m.ScheduledFor,
		Attempts:     m.Attempts,
		MaxAttempts:  m.MaxAttempts,
		LastError:    m.LastError,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}, nil
}

func fromDomain(t *notification.Task) *taskModel {
	dto := metadataDTO{
		ServiceName:   t.Metadata.ServiceName,
		LocationName:  t.Metadata.LocationName,
		AppointmentAt: t.Metadata.AppointmentAt,
	}
	meta, _ := json.Marshal(dto)
	return &taskModel{
		ID:             t.ID,
		AppointmentID:  t.AppointmentID,
		Type:           string(t.Type),
		RecipientType:  string(t.RecipientType),
		RecipientPhone: t.RecipientPhone.String(),
		Metadata:       meta,
		Status:         string(t.Status),
		ScheduledFor:   t.ScheduledFor,
		Attempts:       t.Attempts,
		MaxAttempts:    t.MaxAttempts,
		LastError:      t.LastError,
		CreatedAt:      t.CreatedAt,
		UpdatedAt:      t.UpdatedAt,
	}
}
