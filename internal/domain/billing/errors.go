package billing

import "errors"

var (
	ErrPlanNotFound           = errors.New("plan not found")
	ErrPlanInactive           = errors.New("plan is inactive")
	ErrInvalidPlanCode        = errors.New("invalid plan code")
	ErrPlanNameRequired       = errors.New("plan name is required")
	ErrNegativeLimit          = errors.New("limit cannot be negative")
	ErrPlanCapabilityNotFound = errors.New("plan capability not found")
)

var (
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrSubscriptionActive        = errors.New("subscription is already active")
	ErrSubscriptionExpired       = errors.New("subscription has expired")
	ErrSubscriptionSuspended     = errors.New("subscription is suspended")
	ErrLimitExceeded             = errors.New("plan limit exceeded")
	ErrInvalidStatusTransition   = errors.New("invalid subscription status transition")
	ErrInvalidSubscriptionStatus = errors.New("invalid subscription status")
	ErrMaxRetriesExceeded        = errors.New("maximum payment retry attempts exceeded")
)

var (
	ErrInvalidBillingCycle  = errors.New("invalid billing cycle")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
)
