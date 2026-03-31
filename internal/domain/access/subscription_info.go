package access

import (
	"time"

	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
)

type SubscriptionInfo struct {
	Status           billing.SubscriptionStatus
	CurrentPeriodEnd *time.Time
	LocationsUsed    int
	EmployeesUsed    int
}
