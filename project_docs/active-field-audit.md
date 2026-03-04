# Аудит поля Active — все сущности

Дата: 2026-03-04

## Сводная таблица

| Сущность | Фильтр в SQL | Проверка в коде | Статус |
|----------|-------------|-----------------|--------|
| **Employee** | `active = true` в `getOrgContext` | `IsActive()` в логине, refresh; `CanServeClients()` в booking | OK |
| **Organization** | Нет (by design) | `orgCtx.Organization.Active` в OrgContext middleware | OK |
| **Plan** | `active = true` в `findActiveByCode` | — | OK |
| **Location** | Нет (single-entity load) | `loc.CanOperate()` в booking и transfer | OK |
| **Service** | Нет (single-entity load) | `svc.IsActive()` в booking | OK |
| **EmployeeService** | `active = true` во всех read-запросах | — | OK |
| **ServiceCategory** | Репозиторий не создан | — | Будущая задача |

## Детали по каждой сущности

### Employee

Три точки проверки `Active`:

1. **Login** (`usecases/auth/auth.go`) — `IsActive()` проверяется после валидации пароля. Деактивированный сотрудник не получит токен.
2. **Refresh tokens** (`usecases/auth/auth.go`) — `IsActive()` проверяется перед выдачей новых токенов.
3. **OrgContext middleware** — SQL запрос `getOrgContext` содержит `WHERE active = true AND deleted_at IS NULL`. Даже с валидным токеном (выданным до деактивации) middleware вернёт 401.
4. **Booking (публичный)** — `CanServeClients()` = `Active && !IsDeleted() && CanProvideServices`. Деактивированный сотрудник не появится в слотах и не примет запись.

`CanServeClients()` включает `Active` потому что booking — публичный эндпоинт без auth middleware. Без этой проверки деактивированный сотрудник был бы виден клиентам.

Employee_services **не деактивируются** при деактивации сотрудника — `CanServeClients()` отсекает на уровне выше, а при реактивации все услуги сразу доступны без восстановления.

### Organization

Проверяется в **OrgContext middleware** (`middlewares/org_context.go`). Если `orgCtx.Organization.Active == false` — возвращает 403 `organization_inactive`.

Для публичного booking отдельная проверка не нужна — деактивация организации должна сопровождаться приостановкой подписки (`ErrSubscriptionSuspended`).

### Plan

Фильтруется на уровне SQL: `WHERE active = true` в `findActiveByCode`. Единственный read-запрос. Неактивный план невозможно выбрать при регистрации.

### Location

Проверяется в коде через `loc.CanOperate()` = `Active && !IsDeleted()`:

- **Booking** (`usecases/booking/usecase.go`, `get_available_slots.go`) — нельзя забронировать и получить слоты в неактивной локации (`ErrLocationInactive`).
- **Transfer** (`usecases/employee/transfer_employee.go`) — нельзя перевести сотрудника в неактивную локацию.

SQL не фильтрует `active` — `GetByID` используется и для мутаций (activate/deactivate).

### Service

Проверяется в коде через `svc.IsActive()` = `Active && !IsDeleted()`:

- **Booking** (`usecases/booking/usecase.go`, `get_available_slots.go`) — нельзя забронировать неактивную услугу (`ErrServiceInactive`).

SQL не фильтрует `active` — `GetByID` используется и для мутаций.

### EmployeeService

Фильтруется на уровне SQL: `WHERE active = true` во всех read-запросах (`GetActiveByEmployeeAndService`, `GetActiveByLocationAndService`). Неактивная связка employee-service не возвращается из репозитория.

### ServiceCategory

Репозиторий ещё не создан. При добавлении — list-запросы должны фильтровать `WHERE active = true`.

## Принцип: где фильтровать Active

| Ситуация | Где фильтровать | Почему |
|----------|----------------|--------|
| Read-only запрос (списки, публичные API) | SQL `WHERE active = true` | Неактивные данные не должны покидать БД |
| Single-entity load (`GetByID`) | В коде (use case) | `GetByID` используется и для мутаций (activate/deactivate) |
| Бизнес-логика зависит от нескольких полей | В домене (entity method) | `CanServeClients()`, `CanOperate()` |

## Что было сделано (2026-03-04)

1. Добавлена проверка `IsActive()` в `LoginEmployee()` и `RefreshTokens()`
2. Добавлен `active = true` в SQL `getOrgContext` (middleware)
3. Добавлена проверка `Organization.Active` в OrgContext middleware
4. Добавлена проверка `loc.CanOperate()` в booking use cases
5. Добавлена проверка `svc.IsActive()` в booking use cases
6. `CanServeClients()` = `Active && !IsDeleted() && CanProvideServices` (Active вернулся)
7. Добавлены sentinel errors: `ErrLocationInactive`, `ErrServiceInactive`
8. Error maps обновлены в booking и employee handlers
