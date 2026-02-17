package authrepo

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/auth"
)

type pendingRegistrationModel struct {
	Phone        string    `json:"phone"`
	FirstName    string    `json:"first_name"`
	LastName     *string   `json:"last_name,omitempty"`
	PasswordHash string    `json:"password_hash"`
	PlanCode     string    `json:"plan_code"`
	OTPCode      string    `json:"otp_code"`
	Attempts     int       `json:"attempts"`
	MaxAttempts  int       `json:"max_attempts"`
	LastSentAt   time.Time `json:"last_sent_at"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func fromPendingRegistration(reg *auth.PendingRegistration) *pendingRegistrationModel {
	return &pendingRegistrationModel{
		Phone:        reg.Phone.String(),
		FirstName:    reg.FirstName,
		LastName:     reg.LastName,
		PasswordHash: reg.PasswordHash,
		PlanCode:     reg.PlanCode.Value(),
		OTPCode:      reg.OTPCode,
		Attempts:     reg.Attempts,
		MaxAttempts:  reg.MaxAttempts,
		LastSentAt:   reg.LastSentAt,
		ExpiresAt:    reg.ExpiresAt,
	}
}

func (m *pendingRegistrationModel) toUseCase() (*auth.PendingRegistration, error) {
	phone, err := shared.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}

	planCode, err := shared.NewSlug(m.PlanCode)
	if err != nil {
		return nil, err
	}

	return &auth.PendingRegistration{
		Phone:        phone,
		FirstName:    m.FirstName,
		LastName:     m.LastName,
		PasswordHash: m.PasswordHash,
		PlanCode:     planCode,
		OTPCode:      m.OTPCode,
		Attempts:     m.Attempts,
		MaxAttempts:  m.MaxAttempts,
		LastSentAt:   m.LastSentAt,
		ExpiresAt:    m.ExpiresAt,
	}, nil
}
