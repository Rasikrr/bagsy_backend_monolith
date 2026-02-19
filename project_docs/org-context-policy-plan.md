# План: Внедрение OrgContext и Policy Layer

**Дата:** 18.02.2026
**Статус:** Проектирование / Согласовано

## 1. Цели
1. **OrgContext:** Предоставить UseCase'ам полную "картину" текущего контекста (сотрудник, организация, подписка, план) без лишних походов в БД внутри бизнес-логики.
2. **Policy:** Вынести правила авторизации и проверки лимитов в отдельный тестируемый слой, изолированный от инфраструктуры (репозиториев).

---

## 2. OrgContext (Assembly Strategy)

### Схема сборки (Middleware)
Middleware `OrgContextMiddleware` выполняется строго **после** `AuthMiddleware`.
1. Из `ctx` извлекается `UserID` (через `access.TokenFromContext`).
2. Выполняется **один оптимизированный SQL-запрос** с JOIN: `employees` + `organizations` + `subscriptions` + `plans`.
3. Результат мапится в структуру `access.OrgContext`.
4. Объект сохраняется в `ctx` через `access.WithOrgContext`.

### Репозиторий
Метод для сборки контекста будет находиться в новом репозитории или расширит существующий `employee`:
```go
// Псевдо-запрос
SELECT e.id, e.role, o.id, o.status, s.status, p.code, p.capabilities
FROM employees e
JOIN organizations o ON e.organization_id = o.id
LEFT JOIN subscriptions s ON s.organization_id = o.id AND s.is_active = true
LEFT JOIN plans p ON s.plan_id = p.id
WHERE e.id = $1;
```

---

## 3. Policy Layer (Authorization & Limits)

### Принципы
1. **Stateless (Чистая логика):** Policy не содержит репозиториев. Все данные для проверки она получает извне.
2. **Read-Only:** Policy никогда не меняет состояние системы.
3. **Тестируемость:** Легко покрывается Unit-тестами без моков базы данных.

### Пример реализации
`internal/usecases/policy/policy.go`:
```go
type Policy struct{}

// CanCreateLocation проверяет, может ли организация создать еще одну локацию
func (p *Policy) CanCreateLocation(orgCtx *access.OrgContext, currentLocationsCount int) error {
    // 1. Проверка статуса подписки
    if !orgCtx.Subscription.Status.CanOperate() {
        return billing.ErrSubscriptionSuspended
    }

    // 2. Проверка прав роли
    if !orgCtx.Employee.Role.IsOwner() && !orgCtx.Employee.Role.IsAdmin() {
        return identity.ErrPermissionDenied
    }

    // 3. Проверка лимитов тарифа
    if !orgCtx.Plan.Capabilities.CanUse(billing.ResourceMaxLocations, currentLocationsCount) {
        return billing.ErrLimitExceeded
    }

    return nil
}
```

---

## 4. Взаимодействие в UseCase

UseCase выступает оркестратором: собирает "Usage" (текущее потребление) и передает его в Policy.

```go
func (u *UseCase) CreateLocation(ctx context.Context, params Params) (*Output, error) {
    orgCtx := access.MustFromContext(ctx)
    
    // 1. Загрузка данных для проверки (Usage)
    count, err := u.locationRepo.CountByOrganization(ctx, orgCtx.Organization.ID)
    if err != nil {
        return nil, err
    }
    
    // 2. Валидация через Policy
    if err := u.policy.CanCreateLocation(orgCtx, count); err != nil {
        return nil, err
    }
    
    // 3. Выполнение действия
    loc := location.New(...)
    return u.locationRepo.Save(ctx, loc)
}
```

---

## 5. Порядок реализации

1.  **Repository:** Добавить метод `GetOrgContext` (JOIN запрос).
2.  **Middleware:** Реализовать `OrgContextMiddleware` в `internal/ports/http/middlewares/`.
3.  **Policy:** Создать `internal/usecases/policy/policy.go`.
4.  **DI:** Зарегистрировать новые компоненты в `internal/app/app.go`.
5.  **Unit Tests:** Покрыть `Policy` тестами на разные сценарии (превышение лимита, просроченная подписка).
