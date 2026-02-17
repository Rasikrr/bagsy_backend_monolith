package auth

//go:generate easyjson -all models.go
// ── Register ────────────────────────────────────────────────────

type registerRequest struct {
	Phone     string  `json:"phone"`
	FirstName string  `json:"first_name"`
	LastName  *string `json:"last_name"`
	Password  string  `json:"password"`
	PlanCode  string  `json:"plan_code"`
}

type registerResponse struct {
	Message    string `json:"message"`
	Phone      string `json:"phone"`
	ExpiresIn  int    `json:"expires_in"`
	RetryAfter int    `json:"retry_after"`
}

// ── Verify ──────────────────────────────────────────────────────

type verifyRequest struct {
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type tokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type verifyUser struct {
	ID             string  `json:"id"`
	Phone          string  `json:"phone"`
	FirstName      string  `json:"first_name"`
	LastName       *string `json:"last_name,omitempty"`
	Role           string  `json:"role"`
	OrganizationID string  `json:"organization_id"`
}

// ── Resend ──────────────────────────────────────────────────────

type resendRequest struct {
	Phone string `json:"phone"`
}

type resendResponse struct {
	Message    string `json:"message"`
	ExpiresIn  int    `json:"expires_in"`
	RetryAfter int    `json:"retry_after"`
}

// ── Login ──────────────────────────────────────────────────────

type loginRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
