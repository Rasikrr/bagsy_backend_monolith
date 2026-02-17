# CLAUDE.md — Bagsy Backend Monolith

## Project Overview

Bagsy — SaaS-платформа для управления записями, сотрудниками и услугами организаций (салоны красоты, клиники и т.д.). Multi-tenant модульный монолит на Go.

- **Module path:** `github.com/Rasikrr/bagsy_backend_monolith`
- **Go version:** 1.25.1
- **DB:** PostgreSQL
- **Cache/Sessions:** Redis
- **Object Storage:** S3 (AWS SDK v2)
- **Messaging:** WhatsApp (GreenAPI), SMS

## Architecture

Clean Architecture + DDD. Dependency Rule: зависимости направлены только внутрь.

```
HTTP Request → Handler → UseCase → Domain Entity + Repository → Database
```

### Layer Structure

| Layer | Path | Responsibility |
|-------|------|----------------|
| **Domain** | `internal/domain/{context}/` | Entities, Value Objects, Domain Errors. Zero external deps. |
| **Use Cases** | `internal/usecases/{context}/` | Orchestration: load aggregate → call method → save. |
| **Ports** | `internal/ports/http/` | HTTP handlers, middlewares, server. |
| **Infrastructure** | `internal/infra/` | JWT, external integrations impl. |
| **Packages** | `pkg/` | Shared utilities (hasher, s3, sms, whatsapp). |

### Domain Contexts (Bounded Contexts)

| Context | Package | Key Entities |
|---------|---------|-------------|
| access | `domain/access` | OrgContext (read-only projection for middleware) |
| auth | `domain/auth` | OTP, Token |
| identity | `domain/identity` | Employee, Customer, CustomerBase, Role, Permissions, WorkHistory |
| organization | `domain/organization` | Organization |
| billing | `domain/billing` | Plan, PlanCapability, Subscription, SubscriptionStatus |
| location | `domain/location` | Location, Category, Address, Coordinates |
| schedule | `domain/schedule` | LocationSchedule, EmployeeSchedule, SlotType |
| catalog | `domain/catalog` | Service, Category, EmployeeService |
| booking | `domain/booking` | Appointment, Status, StatusHistory |
| notification | `domain/notification` | Task, Type |
| media | `domain/media` | Asset, Status |
| shared | `domain/shared` | Phone, Money, Slug, Duration (Value Objects) |

## Coding Conventions

### Domain Layer Rules

1. **No framework/infra imports in domain.** Only stdlib + `github.com/google/uuid` + `github.com/shopspring/decimal`.
2. **No `json`, `sql`, `db` tags** on domain structs. Domain is pure Go.
3. **Rich domain models.** Business logic lives in entity methods, not in services.
4. **Value Objects** are immutable structs with private fields and constructor validation (`NewPhone`, `NewMoney`, `NewSlug`).
5. **Domain errors** — `var ErrXxx = errors.New("...")` in `errors.go` per context. Use `errors` stdlib package.
6. **Soft delete pattern** — entities have `DeletedAt *time.Time`. Check `IsDeleted()` before mutations.
7. **`touch()` pattern** — private method sets `UpdatedAt` on every mutation.
8. **Constructor pattern** — `NewXxx(params)` for creation with validation. Use `XxxParams` struct when >3 args.

### Entity Method Structure

```go
// 1. Aggregate struct
type Foo struct { ... }

// 2. Constructor
func NewFoo(params CreateFooParams) (*Foo, error) { ... }

// 3. Business Methods (mutators)
func (f *Foo) DoSomething() error {
    if f.IsDeleted() {
        return ErrFooDeleted
    }
    // business logic
    f.touch()
    return nil
}

// 4. Query Methods (read-only, no touch)
func (f *Foo) IsDeleted() bool { return f.DeletedAt != nil }
func (f *Foo) CanOperate() bool { return f.Active && !f.IsDeleted() }

// 5. Private helpers
func (f *Foo) touch() { now := time.Now(); f.UpdatedAt = &now }
```

### Use Cases Layer Rules

1. One file per use case: `create_appointment.go`, `fire_employee.go`.
2. UseCase struct holds repository/gateway interfaces as dependencies.
3. Always accept `context.Context` as first argument.
4. Orchestration only — load aggregate, call domain method, save. No business rules here.
5. Transactions: wrap multi-entity operations in a single DB transaction.

### Ports (HTTP) Layer Rules

1. Handlers in `internal/ports/http/handlers/{context}/`.
2. Middlewares in `internal/ports/http/middlewares/`.
3. Parse request → validate → call use case → write response.
4. Never return domain entities directly — map to DTOs.
5. Map domain errors to HTTP status codes in handlers.

### Error Handling

Ошибки делятся на два типа: **доменные** (бизнес-логика) и **инфраструктурные** (БД, Redis, сеть, внешние API).

#### 1. Domain Layer — sentinel errors

- Объявляются в `errors.go` каждого bounded context: `var ErrXxx = errors.New("...")`.
- Используется **только** `errors` из stdlib. Никакого `cockroachdb/errors` в domain.
- Sentinel errors представляют бизнес-ситуации: `ErrEmployeeNotFound`, `ErrOTPExpired`, `ErrPhoneAlreadyExists`.

#### 2. Repository / Infra Layer — обёртка через `fmt.Errorf`

- **Инфраструктурные ошибки** (DB/Redis/сеть) оборачиваются через `fmt.Errorf("описание операции: %w", err)`.
- Описание операции — что делал репозиторий: `"save employee"`, `"get refresh token"`, `"find plan by code"`.
- **Никогда** не возвращаем голый `return err` из репозитория — всегда оборачиваем для трассировки.
- Если ошибка означает "не найдено" — возвращаем доменный sentinel: `return nil, identity.ErrEmployeeNotFound`.
- Репозиторий **не импортирует** ошибки внешних библиотек наверх — конвертирует `redis.Nil` → `authDomain.ErrRefreshTokenNotFound`, `pgxscan.NotFound` → `identity.ErrEmployeeNotFound`.

```go
// Repository — правильно
func (r *Repository) GetByPhone(ctx context.Context, phone shared.Phone) (*identity.Employee, error) {
    var m model
    if err := pgxscan.Get(ctx, r.db, &m, getByPhone, phone.String()); err != nil {
        if pgxscan.NotFound(err) {
            return nil, identity.ErrEmployeeNotFound  // доменный sentinel
        }
        return nil, fmt.Errorf("get employee by phone: %w", err)  // инфра обёртка
    }
    return m.toDomain()
}
```

#### 3. Use Case Layer — `cockroachdb/errors` для обёртки

- Оборачивает ошибки инфра-вызовов: `errors.Wrap(err, "описание шага")`.
- **Не** проверяет инфраструктурные ошибки (никогда `errors.Is(err, redis.SomeError)`).
- Проверяет **только** доменные sentinel: `errors.Is(err, authDomain.ErrPhoneAlreadyExists)`.
- Доменные ошибки возвращаются as-is, без обёртки: `return nil, authDomain.ErrOTPAlreadySent`.

#### 4. HTTP Handler Layer — маппинг ошибок

- Каждый handler-пакет имеет свой `errors.go` с `ErrorMap` — маппинг доменных ошибок → HTTP status + клиентский код.
- `internal/ports/http/util/response.go` — общий `SendError()` как fallback (логирует + возвращает 500).
- Клиенту **никогда** не отдаём внутренние сообщения ошибок — только slug-коды: `"otp_expired"`, `"phone_exists"`, `"internal_error"`.
- Полная ошибка (с инфра-деталями) пишется **только** в логи.

```go
// Handler errors.go — декларативная карта
var authErrors = util.ErrorMap{
    authDomain.ErrOTPExpired:         {http.StatusGone, "otp_expired"},
    authDomain.ErrPhoneAlreadyExists: {http.StatusConflict, "phone_exists"},
}

// В хендлере — одна строка
if err != nil {
    util.SendError(ctx, w, err, authErrors)
    return
}
```

#### 5. Pkg Layer (whatsapp, sms, s3) — собственные sentinel errors

- Каждый пакет в `pkg/` объявляет свои sentinel errors в `errors.go`.
- Внутренние ошибки оборачиваются: `errors.Wrap(ErrSendFailed, err.Error())`.
- Пакеты изолированы — не импортируют domain errors.

#### Поток ошибки (сводка)

```
Repository: redis.Nil → authDomain.ErrRefreshTokenNotFound (конвертация)
Repository: pgx timeout → fmt.Errorf("save employee: %w", err) (обёртка)
UseCase:    получает ошибку → errors.Wrap(err, "шаг") или return as-is
Handler:    errors.Is(err, sentinel) → HTTP код + slug
Handler:    неизвестная ошибка → log.Error + 500 "internal_error"
```

### Naming

- **Files:** `snake_case.go`
- **Packages:** singular, lowercase (`identity`, not `identities`)
- **Types:** `PascalCase` — `Employee`, `SubscriptionStatus`
- **Enums:** `type Foo string` with `const FooBar Foo = "bar"` pattern
- **Constructors:** `NewXxx` / `ParseXxx`
- **Boolean methods:** `IsXxx()`, `CanXxx()`, `HasXxx()`, `ShouldXxx()`

### Multi-Tenancy

- `OrgContext` (in `domain/access`) is a read-only projection assembled in middleware.
- All repository queries MUST filter by `organization_id`.
- Authorization checks live in `internal/usecases/policy/`.

### Database / Migrations

- Migrations in `migrations/*.sql` using goose format (`-- +goose Up` / `-- +goose Down`).
- Custom `timerange` type created for `EXCLUDE USING gist` constraints on schedules.
- UUIDs as primary keys (`gen_random_uuid()`).
- Soft deletes via `deleted_at TIMESTAMPTZ`.

### Testing

- Unit tests in `_test.go` files next to source.
- Use `github.com/stretchr/testify` for assertions.
- Domain tests should not require any external dependencies.

## Key Patterns

| Pattern | Where | Purpose |
|---------|-------|---------|
| Aggregate Root | domain entities | Consistency boundary |
| Value Object | `shared/` (Phone, Money, Slug) | Immutable, self-validating |
| Transactional Outbox | notification_outbox table | Reliable event publishing |
| ACL (Anti-Corruption Layer) | `pkg/`, `gateway/` | Isolate external DTOs from domain |
| OrgContext | `domain/access` | Multi-tenant context propagation |
| Policy | `usecases/policy/` | Authorization separated from business logic |

## Commands

```bash
make <target>    # See scripts/*.mk for available targets
```

## Important Files

- `go.mod` — dependencies
- `Makefile` + `scripts/*.mk` — build/migration commands
- `migrations/schema.sql` — full DB schema reference
- `migrations/media.sql` — media assets schema
- `project_docs/` — architecture docs, flow diagrams, mermaid charts
