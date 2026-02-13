package billing

import "errors"

var (
	ErrPlanNotFound     = errors.New("plan not found")
	ErrPlanInactive     = errors.New("plan is inactive")
	ErrInvalidPlanCode  = errors.New("invalid plan code")
	ErrPlanNameRequired = errors.New("plan name is required")
)

var (
	ErrSubscriptionNotFound  = errors.New("subscription not found")
	ErrSubscriptionActive    = errors.New("subscription is already active")
	ErrSubscriptionExpired   = errors.New("subscription has expired")
	ErrSubscriptionSuspended = errors.New("subscription is suspended")
)

var (
	ErrInvalidBillingCycle  = errors.New("invalid billing cycle")
	ErrInvalidPaymentStatus = errors.New("invalid payment status")
)
