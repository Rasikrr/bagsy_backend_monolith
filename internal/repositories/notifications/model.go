package notifications

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/notification"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type model struct {
	ID            uuid.UUID  `db:"id"`
	BagsyID       uuid.UUID  `db:"bagsy_id"`
	Type          string     `db:"type"`
	RecipientType string     `db:"recipient_type"`
	ScheduledAt   time.Time  `db:"scheduled_at"`
	SentAt        *time.Time `db:"sent_at"`
	Status        string     `db:"status"`
	Attempts      int        `db:"attempts"`
	LastError     *string    `db:"last_error"`
	CreatedAt     time.Time  `db:"created_at"`
}

func (m model) convert() (*notification.Notification, error) {
	notifType, err := notification.TypeString(m.Type)
	if err != nil {
		return nil, errors.Wrap(err, "invalid notification type")
	}

	recipientType, err := notification.RecipientTypeString(m.RecipientType)
	if err != nil {
		return nil, errors.Wrap(err, "invalid recipient type")
	}

	status, err := notification.StatusString(m.Status)
	if err != nil {
		return nil, errors.Wrap(err, "invalid notification status")
	}

	return &notification.Notification{
		ID:            m.ID,
		BagsyID:       m.BagsyID,
		Type:          notifType,
		RecipientType: recipientType,
		ScheduledAt:   m.ScheduledAt,
		SentAt:        m.SentAt,
		Status:        status,
		Attempts:      m.Attempts,
		LastError:     m.LastError,
		CreatedAt:     m.CreatedAt,
	}, nil
}

type models []model

func (mm models) convert() ([]*notification.Notification, error) {
	result := make([]*notification.Notification, 0, len(mm))
	for _, m := range mm {
		n, err := m.convert()
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

// nolint: unused
func convertToModel(n *notification.Notification) model {
	return model{
		ID:            n.ID,
		BagsyID:       n.BagsyID,
		Type:          n.Type.String(),
		RecipientType: n.RecipientType.String(),
		ScheduledAt:   n.ScheduledAt,
		SentAt:        n.SentAt,
		Status:        n.Status.String(),
		Attempts:      n.Attempts,
		LastError:     n.LastError,
		CreatedAt:     n.CreatedAt,
	}
}
