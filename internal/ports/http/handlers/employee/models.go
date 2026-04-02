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
	ID           string               `json:"id"`
	Phone        string               `json:"phone"`
	FirstName    string               `json:"first_name"`
	LastName     *string              `json:"last_name,omitempty"`
	AvatarURL    *string              `json:"avatar_url,omitempty"`
	LocationID   *string              `json:"location_id,omitempty"`
	Role         string               `json:"role"`
	Permissions  permissionsResponse  `json:"permissions"`
	Active       bool                 `json:"active"`
	Organization organizationResponse `json:"organization"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    *time.Time           `json:"updated_at,omitempty"`
}

type permissionsResponse struct {
	CanProvideServices        bool `json:"can_provide_services"`
	CanManageLocationSchedule bool `json:"can_manage_location_schedule"`
}

type organizationResponse struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Subscription subscriptionResponse `json:"subscription"`
}

type subscriptionResponse struct {
	Plan             string           `json:"plan"`
	Status           string           `json:"status"`
	CurrentPeriodEnd *time.Time       `json:"current_period_end"`
	Limits           limitsResponse   `json:"limits"`
	Features         featuresResponse `json:"features"`
}

type limitValueResponse struct {
	Used int  `json:"used"`
	Max  *int `json:"max"`
}

type limitsResponse struct {
	Locations       limitValueResponse `json:"locations"`
	Employees       limitValueResponse `json:"employees"`
	BookingsMonthly limitValueResponse `json:"bookings_monthly"`
}

type featuresResponse struct {
	MultiLocation    bool `json:"multi_location"`
	CustomBranding   bool `json:"custom_branding"`
	APIAccess        bool `json:"api_access"`
	SMSNotifications bool `json:"sms_notifications"`
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
	ID                string  `json:"id"`
	EmployeeServiceID string  `json:"employee_service_id"`
	CategoryID        string  `json:"category_id"`
	Name              string  `json:"name"`
	Description       *string `json:"description,omitempty"`
	DurationMinutes   int     `json:"duration_minutes"`
	Color             string  `json:"color"`
	SortOrder         int     `json:"sort_order"`
	Active            bool    `json:"active"`
	Price             *int64  `json:"price,omitempty"`
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
