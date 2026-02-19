package access

import (
	"github.com/Rasikrr/bagsy_backend_monolith/internal/domain/billing"
)

type PlanInfo struct {
	Code         billing.PlanCode
	Capabilities Capabilities
}

type Capabilities struct {
	m map[billing.Resource]billing.Limit // resource → limit
}

func NewCapabilities(m map[billing.Resource]billing.Limit) Capabilities {
	return Capabilities{m: m}
}

// IsAllowed проверяет, разрешена ли конкретная фича (булево).
func (c Capabilities) IsAllowed(feature billing.Resource) bool {
	if c.m == nil {
		return false
	}
	limit, ok := c.m[feature]
	if !ok {
		return false
	}
	// Если лимит безлимитный или больше 0 — разрешено.
	return limit.IsUnlimited() || limit.Value() > 0
}

// CanUse проверяет, можно ли использовать количественный ресурс.
// Возвращает false если ресурса нет в плане или лимит исчерпан.
func (c Capabilities) CanUse(resource billing.Resource, count int) bool {
	if c.m == nil {
		return false
	}
	limit, ok := c.m[resource]
	if !ok {
		return false
	}
	return !limit.IsExceeded(count)
}
