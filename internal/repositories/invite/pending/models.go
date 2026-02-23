package pending

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/identity"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/shared"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/usecases/invite"
	"github.com/google/uuid"
)

type pendingInviteModel struct {
	Phone                     string    `json:"phone"`
	FirstName                 string    `json:"first_name"`
	LastName                  *string   `json:"last_name,omitempty"`
	OrganizationID            string    `json:"organization_id"`
	LocationID                *string   `json:"location_id,omitempty"`
	Role                      string    `json:"role"`
	CanProvideServices        bool      `json:"can_provide_services"`
	CanManageLocationSchedule bool      `json:"can_manage_location_schedule"`
	InvitedBy                 string    `json:"invited_by"`
	LastSentAt                time.Time `json:"last_sent_at"`
	ExpiresAt                 time.Time `json:"expires_at"`
}

func fromPendingInvite(inv *invite.PendingInvite) *pendingInviteModel {
	m := &pendingInviteModel{
		Phone:                     inv.Phone.String(),
		FirstName:                 inv.FirstName,
		LastName:                  inv.LastName,
		OrganizationID:            inv.OrganizationID.String(),
		Role:                      inv.Role.String(),
		CanProvideServices:        inv.Permissions.CanProvideServices,
		CanManageLocationSchedule: inv.Permissions.CanManageLocationSchedule,
		InvitedBy:                 inv.InvitedBy.String(),
		LastSentAt:                inv.LastSentAt,
		ExpiresAt:                 inv.ExpiresAt,
	}
	if inv.LocationID != nil {
		s := inv.LocationID.String()
		m.LocationID = &s
	}
	return m
}

func (m *pendingInviteModel) toUseCase() (*invite.PendingInvite, error) {
	phone, err := shared.NewPhone(m.Phone)
	if err != nil {
		return nil, err
	}
	orgID, err := uuid.Parse(m.OrganizationID)
	if err != nil {
		return nil, err
	}
	role, err := identity.ParseRole(m.Role)
	if err != nil {
		return nil, err
	}
	invitedBy, err := uuid.Parse(m.InvitedBy)
	if err != nil {
		return nil, err
	}

	var locationID *uuid.UUID
	if m.LocationID != nil {
		var id uuid.UUID
		id, err = uuid.Parse(*m.LocationID)
		if err != nil {
			return nil, err
		}
		locationID = &id
	}

	return &invite.PendingInvite{
		Phone:          phone,
		FirstName:      m.FirstName,
		LastName:       m.LastName,
		OrganizationID: orgID,
		LocationID:     locationID,
		Role:           role,
		Permissions:    identity.NewPermissions(m.CanProvideServices, m.CanManageLocationSchedule),
		InvitedBy:      invitedBy,
		LastSentAt:     m.LastSentAt,
		ExpiresAt:      m.ExpiresAt,
	}, nil
}
