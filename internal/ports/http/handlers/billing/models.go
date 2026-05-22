package billing

//go:generate easyjson -all models.go

type activateRequest struct {
	Cycle string `json:"cycle"`
}

type subscriptionResponse struct {
	ID                 string       `json:"id"`
	OrganizationID     string       `json:"organization_id"`
	Status             string       `json:"status"`
	BillingCycle       string       `json:"billing_cycle,omitempty"`
	RecurringAmount    string       `json:"recurring_amount"`
	CurrentPeriodStart *string      `json:"current_period_start,omitempty"`
	CurrentPeriodEnd   *string      `json:"current_period_end,omitempty"`
	NextBillingAt      *string      `json:"next_billing_at,omitempty"`
	CancelAtPeriodEnd  bool         `json:"cancel_at_period_end"`
	RetryCount         int          `json:"retry_count"`
	SuspendedAt        *string      `json:"suspended_at,omitempty"`
	CanceledAt         *string      `json:"canceled_at,omitempty"`
	DataDeleteAt       *string      `json:"data_delete_at,omitempty"`
	CreatedAt          string       `json:"created_at"`
	Plan               planResponse `json:"plan"`
}

type planResponse struct {
	ID           string               `json:"id"`
	Code         string               `json:"code"`
	Name         string               `json:"name"`
	Description  *string              `json:"description,omitempty"`
	PriceMonthly string               `json:"price_monthly"`
	PriceAnnual  string               `json:"price_annual"`
	Capabilities []capabilityResponse `json:"capabilities"`
}

type capabilityResponse struct {
	Resource string `json:"resource"`
	Limit    *int   `json:"limit"`
}

type getSubscriptionResponse struct {
	Subscription subscriptionResponse `json:"subscription"`
}

type listPlansResponse struct {
	Plans []planResponse `json:"plans"`
}
