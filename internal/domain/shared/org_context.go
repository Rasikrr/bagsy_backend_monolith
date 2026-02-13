package shared

import "github.com/google/uuid"

// ═══════════════════════════════════════════════════════════════
// OrgContext — read-only проекция авторизованного сотрудника.
//
// Собирается в middleware из доменных сущностей:
//   - identity.Employee  → EmployeeInfo
//   - organization.Org   → OrganizationInfo
//   - billing.Subscription → SubscriptionInfo (pre-computed)
//   - billing.Plan        → PlanInfo
//
// Не содержит бизнес-логики авторизации.
// Все проверки доступа — в internal/policy/.
// ═══════════════════════════════════════════════════════════════

type OrgContext struct {
	Employee     EmployeeInfo
	Organization OrganizationInfo
	Subscription SubscriptionInfo
	Plan         PlanInfo
}

// ─────────────────────────────────────────────────────────────────
// Role — единственное определение роли сотрудника.
//
// Каноничный источник. identity.Employee хранит shared.Role.
// Дублирование в identity/role.go должно быть удалено.
// ─────────────────────────────────────────────────────────────────

type Role string

const (
	RoleOwner   Role = "owner"
	RoleManager Role = "manager"
	RoleStaff   Role = "staff"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleOwner, RoleManager, RoleStaff:
		return true
	}
	return false
}

func (r Role) String() string {
	return string(r)
}

func ParseRole(s string) (Role, error) {
	role := Role(s)
	if !role.IsValid() {
		return "", ErrInvalidRole
	}
	return role, nil
}

// ─────────────────────────────────────────────────────────────────
// Permissions — гранулярные разрешения сотрудника.
//
// Каноничный источник. identity.Employee хранит shared.Permissions.
// Дублирование в identity/permissions.go должно быть удалено.
// ─────────────────────────────────────────────────────────────────

type Permissions struct {
	CanProvideServices        bool
	CanManageLocationSchedule bool
}

func NewPermissions(canProvideServices, canManageSchedule bool) Permissions {
	return Permissions{
		CanProvideServices:        canProvideServices,
		CanManageLocationSchedule: canManageSchedule,
	}
}

func DefaultPermissions() Permissions {
	return Permissions{}
}

// ─────────────────────────────────────────────────────────────────
// Resource — тип ресурса для plan capabilities.
//
// Тип определяется здесь для типобезопасности PlanInfo.
// Конкретные константы — в billing/resources.go.
// ─────────────────────────────────────────────────────────────────

type Resource string

// ─────────────────────────────────────────────────────────────────
// EmployeeInfo
// ─────────────────────────────────────────────────────────────────

type EmployeeInfo struct {
	ID          uuid.UUID
	LocationID  uuid.UUID
	Role        Role
	Permissions Permissions
}

// ─────────────────────────────────────────────────────────────────
// OrganizationInfo
// ─────────────────────────────────────────────────────────────────

type OrganizationInfo struct {
	ID      uuid.UUID
	OwnerID uuid.UUID
	Active  bool
}

// ─────────────────────────────────────────────────────────────────
// SubscriptionInfo
//
// Все поля pre-computed в middleware из billing.Subscription.
// Логика определения статуса остаётся в billing домене,
// здесь — только результат.
// ─────────────────────────────────────────────────────────────────

type SubscriptionInfo struct {
	Active    bool
	Trialing  bool
	Suspended bool
	Canceled  bool
}

// ─────────────────────────────────────────────────────────────────
// PlanInfo
//
// Capabilities загружаются из billing.Plan + PlanCapability.
// Ключ — shared.Resource для типобезопасности.
// ─────────────────────────────────────────────────────────────────

type PlanInfo struct {
	Code         string
	Capabilities map[Resource]*int // resource → limit (nil = безлимит)
}

// HasCapability проверяет наличие ресурса/фичи в плане.
func (p PlanInfo) HasCapability(resource Resource) bool {
	_, ok := p.Capabilities[resource]
	return ok
}

// GetLimit возвращает лимит ресурса (nil = безлимит).
func (p PlanInfo) GetLimit(resource Resource) *int {
	return p.Capabilities[resource]
}
