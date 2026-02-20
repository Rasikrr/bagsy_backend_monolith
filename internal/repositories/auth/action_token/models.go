package actiontoken

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/auth"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/google/uuid"
)

type model struct {
	Token          string     `json:"token"`
	Phone          string     `json:"phone"`
	LocationID     *uuid.UUID `json:"location_id,omitempty"`
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	Purpose        string     `json:"purpose"`
	ExpiresAt      time.Time  `json:"expires_at"`
}

func (m *model) toDomain() (*auth.ActionToken, error) {
	phone, err := shared.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}

	purpose, err := auth.ParseActionTokenPurpose(m.Purpose)
	if err != nil {
		return nil, err
	}

	return &auth.ActionToken{
		Token:          m.Token,
		Phone:          phone,
		LocationID:     m.LocationID,
		OrganizationID: m.OrganizationID,
		Purpose:        purpose,
		ExpiresAt:      m.ExpiresAt,
	}, nil
}

func fromDomain(token *auth.ActionToken) *model {
	return &model{
		Token:          token.Token,
		Phone:          token.Phone.String(),
		LocationID:     token.LocationID,
		OrganizationID: token.OrganizationID,
		Purpose:        token.Purpose.String(),
		ExpiresAt:      token.ExpiresAt,
	}
}
