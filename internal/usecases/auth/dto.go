package auth

// ── Register ────────────────────────────────────────────────────

type RegisterInput struct {
	Phone     string
	FirstName string
	LastName  *string
	Password  string
	PlanCode  string
}

type RegisterOutput struct {
	Phone      string
	ExpiresIn  int // seconds
	RetryAfter int // seconds
}

// ── Verify ──────────────────────────────────────────────────────

type VerifyInput struct {
	Phone string
	Code  string
}

type TokensOutput struct {
	AccessToken  string
	RefreshToken string
}

// ── Resend ──────────────────────────────────────────────────────

type ResendInput struct {
	Phone string
}

type ResendOutput struct {
	ExpiresIn  int // seconds
	RetryAfter int // seconds
}

// ── Password Reset ────────────────────────────────────────────────

type RequestResetInput struct {
	Phone string
}

type ConfirmResetInput struct {
	Token       string
	NewPassword string
}
