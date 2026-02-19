package employee

//go:generate easyjson -all models.go

import "github.com/google/uuid"

// ── Send Invite ────────────────────────────────────────────────

type sendInviteRequest struct {
	Phone      string     `json:"phone"`
	FirstName  string     `json:"first_name"`
	LastName   *string    `json:"last_name,omitempty"`
	LocationID *uuid.UUID `json:"location_id,omitempty"`
	Role       string     `json:"role"`
}

type sendInviteResponse struct {
	Message   string `json:"message"`
	Phone     string `json:"phone"`
	ExpiresIn int    `json:"expires_in"`
}

// ── Confirm Invite ─────────────────────────────────────────────

type confirmInviteRequest struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

type confirmInviteResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ── Resend Invite ──────────────────────────────────────────────

type resendInviteRequest struct {
	Phone string `json:"phone"`
}

type resendInviteResponse struct {
	Message    string `json:"message"`
	Phone      string `json:"phone"`
	ExpiresIn  int    `json:"expires_in"`
	RetryAfter int    `json:"retry_after"`
}

// ── Verify Token ───────────────────────────────────────────────

type verifyInviteTokenResponse struct {
	Phone     string  `json:"phone"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name,omitempty"`
	Role      string  `json:"role"`
}
