package otp

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type model struct {
	AppointmentID uuid.UUID `json:"appointment_id"`
	Phone         string    `json:"phone"`
	Code          string    `json:"code"`
	ExpiresAt     time.Time `json:"expires_at"`
	Attempts      int       `json:"attempts"`
}

func fromDomain(appointmentID uuid.UUID, o *auth.OTPCode) *model {
	return &model{
		AppointmentID: appointmentID,
		Phone:         o.Phone.String(),
		Code:          o.Code,
		ExpiresAt:     o.ExpiresAt,
		Attempts:      o.Attempts,
	}
}

func (m *model) toDomain() (*auth.OTPCode, error) {
	phone, err := shared.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}

	return &auth.OTPCode{
		Phone:     phone,
		Code:      m.Code,
		ExpiresAt: m.ExpiresAt,
		Attempts:  m.Attempts,
	}, nil
}
