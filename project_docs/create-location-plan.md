# План: Создание локации + проверка лимитов плана

**Дата:** 17.02.2026
**Контекст:** После регистрации owner должен создать первую локацию. Перед созданием проверяются ограничения тарифного плана.

---

## Текущее состояние

### Готово

| Компонент | Путь | Что делает |
|-----------|------|-----------|
| `access.OrgContext` | `domain/access/org_context.go` | Структура с Employee, Organization, Subscription, Plan |
| `access.Capabilities` | `domain/access/plan_info.go` | `CanUse(resource, count)`, `IsAllowed(feature)` |
| `access.WithToken` / `TokenFromContext` | `domain/access/context.go` | Кладём/достаём auth.Token из ctx |
| `access.WithOrgContext` / `OrgContextFromContext` | `domain/access/context.go` | Кладём/достаём OrgContext из ctx |
| `billing.Resource` + `Limit` | `domain/billing/` | `ResourceMaxLocations`, `Limit.IsExceeded(count)` |
| `SubscriptionStatus.CanOperate()` | `domain/billing/subscription_status.go` | trial/active/past_due → true |
| `location.NewLocation()` | `domain/location/location.go` | Конструктор с валидацией |
| Auth middleware | `ports/http/middlewares/employees.go` | JWT → auth.Token → ctx |
| Error mapping (auth) | `ports/http/handlers/auth/errors.go` | `ErrorMap` для auth handlers |

### Нужно реализовать

Порядок — снизу вверх (domain → infra → use case → HTTP).

---

## Шаг 1: `billing/errors.go` — sentinel error

Добавить в `internal/domain/billing/errors.go`:

```go
var ErrLimitExceeded = errors.New("plan limit exceeded")
```

---

## Шаг 2: `billing/resources.go` — недостающие ресурсы

Добавить константы для ресурсов из seed data, которые понадобятся в ближайшее время:

```go
ResourceMaxServices Resource = "max_services"
```

Остальные (`online_booking`, `whatsapp_notifications`, `client_base`, `analytics_basic/advanced`, `multi_location_management`) — по мере реализации фич.

---

## Шаг 3: Location repository

Создать `internal/repositories/location/`:

```
internal/repositories/location/
├── repository.go    — Repository struct + методы
├── models.go        — DB model + toDomain/fromDomain
└── statements.go    — SQL-запросы
```

Методы:

```go
type Repository struct { db *postgres.Postgres }

func (r *Repository) Save(ctx context.Context, loc *location.Location) error
func (r *Repository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int, error)
```

`CountByOrganization` — нужен для policy (проверка лимита). Считает только `deleted_at IS NULL`.

---

## Шаг 4: `usecases/policy/` — проверки доступа

Создать `internal/usecases/policy/policy.go`:

```go
type Policy struct {
    locationRepo locationRepository
    // позже: employeeRepo, serviceRepo и т.д.
}

func (p *Policy) CanCreateLocation(ctx context.Context, orgCtx *access.OrgContext) error
```

Логика `CanCreateLocation`:

1. `orgCtx.Subscription.Status.CanOperate()` → иначе `billing.ErrSubscriptionSuspended`
2. `orgCtx.Employee.Role` → проверка что owner или admin (пока только owner)
3. `locationRepo.CountByOrganization(ctx, orgCtx.Organization.ID)` → текущее количество
4. `orgCtx.Plan.Capabilities.CanUse(billing.ResourceMaxLocations, count)` → иначе `billing.ErrLimitExceeded`

---

## Шаг 5: OrgContext middleware

Создать `internal/ports/http/middlewares/org_context.go`.

Этот middleware ставится **после** Auth middleware. По `auth.Token.UserID` загружает из БД:

1. Employee (по ID) → `EmployeeInfo`
2. Organization (по `employee.OrganizationID`) → `OrganizationInfo`
3. Subscription (по `organization.ID`, активная) → `SubscriptionInfo`
4. Plan (по `subscription.PlanID` + capabilities) → `PlanInfo`

Собирает `access.OrgContext` и кладёт в ctx через `access.WithOrgContext`.

Нужные repository-интерфейсы:

```go
type employeeRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*identity.Employee, error)
}

type organizationRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*organization.Organization, error)
}

type subscriptionRepository interface {
    GetActiveByOrg(ctx context.Context, orgID uuid.UUID) (*billing.Subscription, error)
}

type planRepository interface {
    GetByID(ctx context.Context, id uuid.UUID) (*billing.Plan, error)
}
```

> **Замечание:** часть этих методов в репозиториях ещё не реализована (GetByID, GetActiveByOrg). Нужно дописать.

---

## Шаг 6: Location use case

Создать `internal/usecases/location/create_location.go`:

```go
type UseCase struct {
    locationRepo locationRepository
    policy       *policy.Policy
}

func (u *UseCase) CreateLocation(ctx context.Context, orgCtx *access.OrgContext, params CreateLocationInput) (*CreateLocationOutput, error) {
    // 1. Policy check
    if err := u.policy.CanCreateLocation(ctx, orgCtx); err != nil {
        return nil, err
    }

    // 2. Domain: создание локации
    loc, err := location.NewLocation(location.CreateLocationParams{
        OrganizationID: orgCtx.Organization.ID,
        ...
    })

    // 3. Save
    if err := u.locationRepo.Save(ctx, loc); err != nil {
        return nil, errors.Wrap(err, "save location")
    }

    return &CreateLocationOutput{ID: loc.ID}, nil
}
```

---

## Шаг 7: Location HTTP handler

Создать `internal/ports/http/handlers/location/`:

```
internal/ports/http/handlers/location/
├── handler.go       — Handler struct + Init(router)
├── errors.go        — ErrorMap для location
├── models.go        — request/response DTOs
└── create.go        — POST /api/v1/locations
```

`errors.go`:

```go
var locationErrors = util.ErrorMap{
    billing.ErrLimitExceeded:        {Code: http.StatusForbidden, Message: "limit_exceeded"},
    billing.ErrSubscriptionSuspended: {Code: http.StatusForbidden, Message: "subscription_suspended"},
    location.ErrNameRequired:         {Code: http.StatusBadRequest, Message: "name_required"},
    location.ErrInvalidScheduleType:  {Code: http.StatusBadRequest, Message: "invalid_schedule_type"},
    shared.ErrInvalidPhone:           {Code: http.StatusBadRequest, Message: "invalid_phone"},
    // ...
}
```

Handler:

```go
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    orgCtx, ok := access.OrgContextFromContext(ctx)
    if !ok {
        util.SendError(ctx, w, ErrUnauthorized, locationErrors)
        return
    }

    var req createRequest
    if err := coreHTTP.GetData(r, &req); err != nil {
        util.SendBadRequest(ctx, w, err)
        return
    }

    out, err := h.useCase.CreateLocation(ctx, orgCtx, ...)
    if err != nil {
        util.SendError(ctx, w, err, locationErrors)
        return
    }

    coreHTTP.SendData(ctx, w, out, http.StatusCreated)
}
```

Роутинг:

```go
router.Route("/api/v1/locations", func(r chi.Router) {
    r.Use(authMiddleware.Handle)
    r.Use(orgContextMiddleware.Handle)
    r.Post("/", h.create)
})
```

---

## Шаг 8: Недостающие repository-методы

Дописать в существующих репозиториях:

| Репозиторий | Метод | Для чего |
|-------------|-------|----------|
| `employee` | `GetByID(ctx, id)` | OrgContext middleware |
| `organization` | `GetByID(ctx, id)` | OrgContext middleware |
| `subscription` | `GetActiveByOrg(ctx, orgID)` | OrgContext middleware |
| `plan` | `GetByID(ctx, id)` | OrgContext middleware |

---

## Шаг 9: Регистрация в DI (`app.go`)

Подключить в `internal/app/app.go`:

- Location repository
- Policy
- Location use case
- Location handler
- OrgContext middleware
- Роутинг с middleware chain

---

## Порядок реализации (рекомендуемый)

```
1. billing/errors.go         — ErrLimitExceeded (1 строка)
2. billing/resources.go      — ResourceMaxServices (1 строка)
3. Repository methods         — GetByID для employee, org, sub, plan
4. Location repository        — Save, CountByOrganization
5. OrgContext middleware       — сборка OrgContext из БД
6. Policy                     — CanCreateLocation
7. Location use case          — CreateLocation
8. Location handler + errors  — HTTP layer
9. DI wiring в app.go         — собираем всё вместе
```

Шаги 1-2 — минуты. Шаги 3-5 — основная работа. Шаги 6-9 — straightforward.
