# Calendar Feature — Migration Plan (main -> full-refactor)

## 1. Что делает календарь на main

**Эндпоинт:** `GET /api/v1/calendar` (authenticated)

**Назначение:** Возвращает список записей (bagsy) за указанный период с информацией об услуге. Используется сотрудниками для просмотра расписания.

### Флоу на main

```
Handler (getCalendar)
  → parse query params (from, to, point_code?, master_phone?)
  → Service.GetCalendar(ctx, query)
      → validate date range (max 35 дней)
      → buildCalendarFilter (role-based scoping):
          Staff    → только свои записи (filter by own phone)
          Manager  → записи своей точки (optional filter by master)
          Owner    → записи любой точки в сети (filter by point_code + optional master)
      → bagsiesRepository.GetOccupiedSlots(ctx, filter)
      → servicesService.GetByIDs(ctx, serviceIDs)  -- обогащение данными услуг
      → merge → []CalendarElement{ Bagsy, Service }
  → map to DTO → response
```

### Фильтр (OccupiedSlotsFilter на main)

```go
type OccupiedSlotsFilter struct {
    PointCode    string
    MasterPhones []string
    StartAt      time.Time
    EndAt        time.Time
}
```

### Ответ (CalendarElement)

Каждый элемент содержит:
- `bagsy_info`: id, point_code, client_phone, master_phone, status, price, start_at, end_at, comment, timestamps
- `service_info`: id, name, description, duration_minutes, color

---

## 2. Проблемы и замечания по старой реализации

### 2.1. Архитектурные

| # | Проблема | Описание |
|---|----------|----------|
| 1 | **Нет use case слоя** | Handler вызывает Service напрямую. В новой архитектуре нужен отдельный use case. |
| 2 | **Доменные сущности текут в ответ** | `newBagsyInfoDTO(b *bagsy.Bagsy)` напрямую маппит domain entity → JSON DTO в handler. В новой архитектуре DTO создаются в handler из use case output. |
| 3 | **Денормализованные phone** | Старый Bagsy хранит `ClientPhone`/`MasterPhone` строками. Новый Appointment хранит `CustomerID`/`EmployeeID` (UUID). Для отображения имени/телефона потребуется JOIN или дополнительные запросы. |

### 2.2. Бизнес-логика

| # | Проблема | Описание |
|---|----------|----------|
| 4 | **Нет фильтрации по статусу** | Репозиторий возвращает ВСЕ записи за период (включая cancelled). Для календаря cancelled записи скорее всего не нужны, либо нужен опциональный фильтр. |
| 5 | **Нет пагинации** | При большом количестве записей ответ может быть огромным. Для MVP допустимо (35 дней × N мастеров), но стоит помнить. |
| 6 | **point_code → locationID** | Старая модель использует `point_code` (строку). Новая — `locationID` (UUID). Фильтрация меняется. |
| 7 | **network_code scoping отсутствует** | Owner мог видеть любую точку в сети, проверяя `point.NetworkCode == act.NetworkCode()`. В новой модели scoping идёт через `organization_id`, что проще и надёжнее. |

### 2.3. Оптимизации

| # | Проблема | Рекомендация |
|---|----------|-------------|
| 8 | **N+1 запрос на services** | Старый код собирает serviceIDs, потом делает `GetByIDs`. Лучше сделать один SQL JOIN. |
| 9 | **N+1 на employee/customer info** | Для отображения имён нужны доп. запросы. Оптимально — JOIN в одном SQL. |

---

## 3. Маппинг: main → full-refactor

| main (старое) | full-refactor (новое) |
|---|---|
| `bagsy.Bagsy` | `booking.Appointment` |
| `bagsy.Status` (uint8 enum) | `booking.Status` (string enum) |
| `service.Service` | `catalog.Service` |
| `point.Point` / `PointCode` | `location.Location` / `LocationID` |
| `user.User` / `MasterPhone` | `identity.Employee` / `EmployeeID` |
| `user.User` / `ClientPhone` | `identity.Customer` / `CustomerID` |
| `actor.Actor` (context) | `access.OrgContext` (context) |
| `actor.Role` (Staff/Manager/SelfOwner/NetManager) | `identity.Role` (staff/manager/owner) |
| `bagsies.Service` | `booking.UseCase` |
| `bagsies.CalendarElement` | Use case output DTO |

---

## 4. План реализации

### Шаг 1: Use Case DTO (`internal/usecases/booking/dto.go`)

Добавить input/output структуры:

```go
type GetCalendarInput struct {
    OrganizationID uuid.UUID
    LocationID     *uuid.UUID  // nil = все локации организации (для owner)
    EmployeeID     *uuid.UUID  // nil = все сотрудники
    StartDate      time.Time
    EndDate        time.Time
    IncludeCancelled bool       // опционально, default false
}

type CalendarEntry struct {
    // Appointment data
    AppointmentID   uuid.UUID
    Status          string
    StartAt         time.Time
    EndAt           time.Time
    Price           float64
    DurationMinutes int
    CustomerComment *string

    // Employee data (denormalized for response)
    EmployeeID   uuid.UUID
    EmployeeName string

    // Customer data
    CustomerID    uuid.UUID
    CustomerName  string
    CustomerPhone string

    // Service data
    ServiceID   uuid.UUID
    ServiceName string
    ServiceColor string

    // Location data
    LocationID   uuid.UUID
    LocationName string
}

type GetCalendarOutput struct {
    Entries []CalendarEntry
}
```

### Шаг 2: Repository — новый метод (`internal/repositories/booking/`)

Добавить SQL-запрос с JOIN-ами для получения всех данных за один запрос:

**`statements.go`** — новый запрос `getCalendarEntries`:

```sql
SELECT
    a.id, a.status, a.start_at, a.end_at, a.price, a.duration_minutes,
    a.customer_comment, a.location_id,
    a.employee_id, e.first_name || COALESCE(' ' || e.last_name, '') AS employee_name,
    a.customer_id, c.first_name || COALESCE(' ' || c.last_name, '') AS customer_name,
    c.phone AS customer_phone,
    a.service_id, s.name AS service_name, s.color AS service_color,
    l.name AS location_name
FROM appointments a
JOIN employees e ON e.id = a.employee_id
JOIN customers c ON c.id = a.customer_id
JOIN services s ON s.id = a.service_id
JOIN locations l ON l.id = a.location_id
WHERE a.organization_id = $1
  AND a.start_at < $3
  AND a.end_at > $2
  AND ($4::uuid IS NULL OR a.location_id = $4)
  AND ($5::uuid IS NULL OR a.employee_id = $5)
  AND ($6::boolean IS TRUE OR a.status != 'cancelled')
ORDER BY a.start_at ASC
```

**`repository.go`** — новый метод:

```go
func (r *Repository) GetCalendarEntries(ctx context.Context, params CalendarQueryParams) ([]CalendarEntryRow, error)
```

**`models.go`** — добавить `calendarEntryRow` модель с `db` тегами для scan.

### Шаг 3: Use Case — метод GetCalendar (`internal/usecases/booking/get_calendar.go`)

Новый файл:

```go
func (u *UseCase) GetCalendar(ctx context.Context, orgCtx *access.OrgContext, input GetCalendarInput) (*GetCalendarOutput, error) {
    // 1. Validate date range
    if input.EndDate.Before(input.StartDate) {
        return nil, booking.ErrInvalidTimeRange
    }
    days := int(input.EndDate.Sub(input.StartDate).Hours()/24) + 1
    if days > maxCalendarRangeDays {
        return nil, booking.ErrCalendarRangeTooLarge  // новая доменная ошибка
    }

    // 2. Apply role-based scoping (через policy)
    //    Staff  → принудительно EmployeeID = свой ID
    //    Manager → принудительно LocationID = свой LocationID
    //    Owner  → без ограничений в рамках org

    // 3. Query repository (один запрос с JOIN)

    // 4. Map to output
}
```

**Scoping правила:**

| Роль | LocationID | EmployeeID |
|------|-----------|------------|
| `staff` | = свой LocationID | = свой EmployeeID |
| `manager` | = свой LocationID | из input (опционально) |
| `owner` | из input (опционально) | из input (опционально) |

### Шаг 4: Policy (`internal/usecases/policy/`)

Добавить метод:

```go
func (p *Policy) CanViewCalendar(orgCtx *access.OrgContext, input *GetCalendarInput) error
```

Или вместо отдельного policy-метода — применить scoping прямо в use case (проще для MVP, scoping это не совсем authorization а скорее data filtering).

**Рекомендация:** scoping правила применять в use case, не в policy. Policy отвечает за "можно или нельзя", а scoping — за "что именно видишь". Все роли МОГУТ смотреть календарь, разница только в scope.

### Шаг 5: Доменная ошибка

**`internal/domain/booking/errors.go`** — добавить:

```go
var ErrCalendarRangeTooLarge = errors.New("calendar range exceeds maximum allowed days")
```

### Шаг 6: Handler (`internal/ports/http/handlers/booking/`)

**`get_calendar.go`** — новый файл:

```go
// GET /api/v1/bookings/calendar?from=2026-01-01&to=2026-01-31&location_id=...&employee_id=...
func (h *Handler) getCalendar(w http.ResponseWriter, r *http.Request)
```

**`models.go`** — добавить request/response модели для calendar.

**`handler.go`** (router) — добавить route в authenticated group:

```go
admin.Get("/calendar", h.getCalendar)
```

**`errors.go`** — добавить маппинг:

```go
booking.ErrCalendarRangeTooLarge: {Code: http.StatusBadRequest, Message: "calendar_range_too_large"},
```

### Шаг 7: Wiring (`internal/app/app.go`)

Никаких новых зависимостей не нужно — booking use case уже имеет `appointmentRepo`, а новый метод `GetCalendarEntries` будет на том же репозитории. Нужно только добавить метод в интерфейс `appointmentRepository` в use case.

---

## 5. Файлы к изменению/созданию

| Файл | Действие |
|------|----------|
| `internal/domain/booking/errors.go` | добавить `ErrCalendarRangeTooLarge` |
| `internal/usecases/booking/dto.go` | добавить `GetCalendarInput`, `CalendarEntry`, `GetCalendarOutput` |
| `internal/usecases/booking/get_calendar.go` | **создать** — use case метод |
| `internal/usecases/booking/usecase.go` | расширить интерфейс `appointmentRepository` |
| `internal/repositories/booking/statements.go` | добавить SQL `getCalendarEntries` |
| `internal/repositories/booking/models.go` | добавить `calendarEntryRow` |
| `internal/repositories/booking/repository.go` | добавить метод `GetCalendarEntries` |
| `internal/ports/http/handlers/booking/get_calendar.go` | **создать** — handler |
| `internal/ports/http/handlers/booking/models.go` | добавить request/response DTO |
| `internal/ports/http/handlers/booking/handler.go` | добавить route |
| `internal/ports/http/handlers/booking/errors.go` | добавить маппинг ошибки |

---

## 6. Что НЕ переносим (отличия от main)

1. **`point_code` фильтрацию** — заменяется на `location_id` (UUID).
2. **`network_code` проверку** — не нужна, scoping через `organization_id`.
3. **`easyjson` генерацию** — на full-refactor не используется.
4. **Отдельный `calendarService`** — календарь становится методом на `booking.UseCase`.
5. **`CalendarElement` с domain entities** — заменяется на плоский DTO из одного SQL.

---

## 7. Verification

```bash
go build ./...
go test ./internal/usecases/booking/...
go vet ./...
```

Ручное тестирование:
- `GET /api/v1/bookings/calendar?from=2026-02-01&to=2026-02-28` с токеном staff — видит только свои записи
- Тот же запрос с токеном manager — видит записи своей локации
- Тот же запрос с токеном owner + `location_id=...` — видит записи конкретной локации
- Запрос с диапазоном > 35 дней — 400
- Запрос с `from > to` — 400
