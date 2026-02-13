package access

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
// Все проверки доступа — в internal/usecases/policy/.
// ═══════════════════════════════════════════════════════════════

type OrgContext struct {
	Employee     EmployeeInfo
	Organization OrganizationInfo
	Subscription SubscriptionInfo
	Plan         PlanInfo
}
