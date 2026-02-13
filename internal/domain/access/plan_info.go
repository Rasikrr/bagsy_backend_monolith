package access

type PlanInfo struct {
	Code         string
	Capabilities Capabilities
}

// ─────────────────────────────────────────────────────────────────
// Capabilities — набор ресурсов/фич плана с лимитами.
//
// Value object. Nil-safe: безопасно вызывать методы
// даже если внутренняя map не инициализирована.
// Ключи — строковые константы из billing.Resource.
// ─────────────────────────────────────────────────────────────────

type Capabilities struct {
	m map[string]*int // resource → limit (nil = безлимит)
}

func NewCapabilities(m map[string]*int) Capabilities {
	return Capabilities{m: m}
}

// Has проверяет наличие ресурса/фичи в плане.
func (c Capabilities) Has(resource string) bool {
	if c.m == nil {
		return false
	}
	_, ok := c.m[resource]
	return ok
}

// GetLimit возвращает лимит ресурса (nil = безлимит).
func (c Capabilities) GetLimit(resource string) *int {
	if c.m == nil {
		return nil
	}
	return c.m[resource]
}

// CheckLimit проверяет, не превышен ли лимит.
// Возвращает true если ресурс безлимитный или count < limit.
// Возвращает false если ресурса нет в плане.
func (c Capabilities) CheckLimit(resource string, count int) bool {
	if c.m == nil {
		return false
	}
	limit, ok := c.m[resource]
	if !ok {
		return false
	}
	if limit == nil {
		return true
	}
	return count < *limit
}
