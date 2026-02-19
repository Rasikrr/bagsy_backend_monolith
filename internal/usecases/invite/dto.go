package invite

import "github.com/google/uuid"

// ── Send Invite ────────────────────────────────────────────────

type SendInviteInput struct {
	Phone      string
	FirstName  string
	LastName   *string
	LocationID *uuid.UUID
	Role       string
}

type SendInviteOutput struct {
	Phone     string
	ExpiresIn int // seconds
}

// ── Confirm Invite ─────────────────────────────────────────────

type ConfirmInviteInput struct {
	Token    string
	Password string
}

type TokensOutput struct {
	AccessToken  string
	RefreshToken string
}

// ── Resend Invite ──────────────────────────────────────────────

type ResendInviteInput struct {
	Phone string
}

type ResendInviteOutput struct {
	Phone      string
	ExpiresIn  int // seconds
	RetryAfter int // seconds
}

// ── Verify Token ───────────────────────────────────────────────

type VerifyInviteTokenOutput struct {
	Phone     string
	FirstName string
	LastName  *string
	Role      string
}
