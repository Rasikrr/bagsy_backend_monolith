package employee

//go:generate easyjson -all models.go

import (
	"time"

	"github.com/google/uuid"
)

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

// ── Get Me ───────────────────────────────────────────────────────

type getMeResponse struct {
	ID             string              `json:"id"`
	Phone          string              `json:"phone"`
	FirstName      string              `json:"first_name"`
	LastName       *string             `json:"last_name,omitempty"`
	AvatarURL      *string             `json:"avatar_url,omitempty"`
	OrganizationID string              `json:"organization_id"`
	LocationID     *string             `json:"location_id,omitempty"`
	Role           string              `json:"role"`
	Permissions    permissionsResponse `json:"permissions"`
	Active         bool                `json:"active"`
	CreatedAt      time.Time           `json:"created_at"`
}

type permissionsResponse struct {
	CanProvideServices        bool `json:"can_provide_services"`
	CanManageLocationSchedule bool `json:"can_manage_location_schedule"`
}

// ── Update Me ───────────────────────────────────────────────────

type updateMeRequest struct {
	FirstName string     `json:"first_name"`
	LastName  *string    `json:"last_name,omitempty"`
	AvatarID  *uuid.UUID `json:"avatar_id,omitempty"`
}

// ── Admin Actions ───────────────────────────────────────────────

type transferRequest struct {
	LocationID uuid.UUID `json:"location_id"`
}

type changeRoleRequest struct {
	Role string `json:"role"`
}

type changePermissionsRequest struct {
	CanProvideServices        bool `json:"can_provide_services"`
	CanManageLocationSchedule bool `json:"can_manage_location_schedule"`
}

// ── Get Employee Services ─────────────────────────────────────────

type employeeServiceItemResponse struct {
	ID              string  `json:"id"`
	CategoryID      string  `json:"category_id"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	DurationMinutes int     `json:"duration_minutes"`
	Color           string  `json:"color"`
	SortOrder       int     `json:"sort_order"`
	Active          bool    `json:"active"`
	Price           *int64  `json:"price,omitempty"`
}

type getEmployeeServicesResponse struct {
	Services []employeeServiceItemResponse `json:"services"`
}

// ── Get List ──────────────────────────────────────────────────────

type employeeListItemResponse struct {
	ID          string              `json:"id"`
	Phone       string              `json:"phone"`
	FirstName   string              `json:"first_name"`
	LastName    *string             `json:"last_name,omitempty"`
	AvatarURL   *string             `json:"avatar_url,omitempty"`
	LocationID  *string             `json:"location_id,omitempty"`
	Role        string              `json:"role"`
	Permissions permissionsResponse `json:"permissions"`
	Active      bool                `json:"active"`
	CreatedAt   time.Time           `json:"created_at"`
}

type getListResponse struct {
	Employees []employeeListItemResponse `json:"employees"`
	Total     int                        `json:"total"`
}
