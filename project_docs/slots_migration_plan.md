# Миграция эндпоинтов Slots на ветку full-refactor

## 1. Анализ текущей реализации (main)

### Эндпоинты

| Эндпоинт | Назначение |
|-----------|-----------|
| `POST /api/v1/bagsies/slots` | Возвращает список доступных **дат** на 28 дней вперёд |
| `POST /api/v1/bagsies/slots/day` | Возвращает **тайм-слоты** по конкретной дате, сгруппированные по мастерам |

### Алгоритм генерации слотов (main)

1. Загрузить Service (для `duration_minutes`), Point (для расписания точки), мастеров (MasterService + User)
2. Из БД получить занятые баgsies (не удалённые, не отменённые) с перекрытием по запрошенному диапазону
3. Для каждого мастера, для каждого дня:
   - Найти пересечение расписания точки и расписания мастера → эффективные рабочие часы
   - Генерировать слоты с шагом 30 мин, длительностью `duration_minutes`
   - Отфильтровать прошедшие и перекрывающиеся с занятыми записями
4. Handler `/slots` — собирает уникальные даты из всех слотов, конвертирует в Almaty TZ
5. Handler `/slots/day` — группирует слоты по мастерам, форматирует время в `"15:04"`

### Ключевые зависимости (main)

- `bagsy.Bagsy` — сущность записи (booking)
- `point.Point` + `point.Schedule` — расписание точки (JSONB, по дням недели)
- `user.User` + `user.Schedule` — расписание мастера (JSONB, по дням недели)
- `masterservice.MasterService` — связь мастер-услуга с ценой
- `service.Service` — длительность услуги
- Идентификация точки по `point_code` (string slug)

---

## 2. Что уже есть на full-refactor

### Готово

| Компонент | Статус | Файл |
|-----------|--------|------|
| Domain entity `Appointment` | Готов | `internal/domain/booking/appointment.go` |
| Domain entity `Location` | Готов | `internal/domain/location/location.go` |
| Domain entity `Employee` | Готов | `internal/domain/identity/employee.go` |
| Domain entity `Service` | Готов | `internal/domain/catalog/service.go` |
| Domain entity `EmployeeService` | Готов | `internal/domain/catalog/employee_service.go` |
| Domain `LocationScheduleSlot` / `EmployeeScheduleSlot` | Готов | `internal/domain/schedule/` |
| UseCase DTO: `GetAvailableSlotsInput/Output`, `TimeSlot` | Готов | `internal/usecases/booking/dto.go` |
| Алгоритм `generateSlots()` | Готов | `internal/usecases/booking/slots.go` |
| Handler DTO: `getSlotsRequest/getSlotsResponse` | Готов | `internal/ports/http/handlers/booking/models.go` |
| Repository интерфейсы (все нужные методы объявлены) | Готов | `internal/usecases/booking/usecase.go` |

### Не готово

| Компонент | Статус | Что нужно |
|-----------|--------|-----------|
| UseCase метод `GetAvailableSlots` | Отсутствует | Оркестрация: загрузка данных → вызов `generateSlots()` → сборка output |
| Handler метод для slots | Отсутствует | Новый файл `get_slots.go` |
| Регистрация роута | Отсутствует | Добавить `POST /api/v1/bookings/slots` в `handler.go` |
| Interface `bookingUseCase` в handler | Не содержит метод для slots | Добавить `GetAvailableSlots` |
| Infra-реализации репозиториев | Отсутствуют все | PostgreSQL/Redis имплементации |
| Error mapping для slots | Частично | Возможно нужны дополнительные маппинги |

---

## 3. Ключевые архитектурные изменения (main → full-refactor)

| Аспект | main (старый) | full-refactor (новый) |
|--------|---------------|----------------------|
| Идентификация точки | `point_code` (string slug) | `location_id` (UUID) |
| Сущность записи | `Bagsy` | `Appointment` |
| Мастер | `User` по `phone` | `Employee` по `UUID` |
| Расписание | JSONB по дням недели в Point/User | Отдельный домен `schedule` с `LocationScheduleSlot` / `EmployeeScheduleSlot` (по конкретным датам, с типами work/rest) |
| Тип расписания | Всегда пересечение point + master | `ScheduleType`: `fixed` (только Location) / `mixed` (пересечение Location + Employee) |
| Шаг слота | Константа 30 мин | `Location.SlotDurationMinutes` (настраиваемый) |
| API shape | 2 эндпоинта (dates + day) | 1 универсальный эндпоинт с диапазоном дат |
| Цена | `decimal.Decimal` напрямую | `shared.Money` (value object с валютой) |

---

## 4. План реализации

### Шаг 1: UseCase метод `GetAvailableSlots`

**Файл:** `internal/usecases/booking/get_available_slots.go`

Оркестрация:
1. `locationRepo.GetByID()` → получить `Location` (`ScheduleType`, `SlotDurationMinutes`)
2. `serviceRepo.GetByID()` → получить `Service` (`DurationMinutes`)
3. Определить сотрудников:
   - Если `input.EmployeeID != nil` → использовать одного
   - Иначе → `empServiceRepo.GetByLocationAndService(locationID, serviceID)` → список сотрудников
4. `employeeRepo.GetByIDs()` → загрузить данные сотрудников (имена)
5. `scheduleRepo.GetLocationSlots(locationID, startDate, endDate)` → расписание локации
6. `scheduleRepo.GetEmployeesSlots(employeeIDs, startDate, endDate)` → расписания сотрудников
7. `appointmentRepo.GetOccupiedSlots(locationID, employeeIDs, startDate, endDate)` → занятые слоты
8. Для каждого сотрудника вызвать `generateSlots()` с соответствующими данными
9. Собрать `GetAvailableSlotsOutput`

### Шаг 2: Handler

**Файл:** `internal/ports/http/handlers/booking/get_slots.go`

- Парсинг `getSlotsRequest` (DTO уже определён)
- Маппинг в `GetAvailableSlotsInput`
- Вызов `useCase.GetAvailableSlots()`
- Маппинг результата в `getSlotsResponse`
- Swagger-аннотации

### Шаг 3: Регистрация роута и интерфейса

**Файл:** `internal/ports/http/handlers/booking/handler.go`

- Добавить метод `GetAvailableSlots` в интерфейс `bookingUseCase`
- Зарегистрировать `POST /api/v1/bookings/slots` как публичный роут (без auth middleware)

### Шаг 4: Error mapping

**Файл:** `internal/ports/http/handlers/booking/errors.go`

Добавить маппинги для возможных ошибок:
- `location.ErrLocationNotFound` → 404
- `catalog.ErrServiceNotFound` → 404
- `identity.ErrEmployeeNotFound` → 404
- `catalog.ErrEmployeeServiceNotFound` → 404

### Шаг 5: Infra-реализации репозиториев (отдельная задача)

Для работы slots нужны реализации следующих интерфейсов:

| Интерфейс | Нужные методы для slots |
|-----------|------------------------|
| `appointmentRepository` | `GetOccupiedSlots(ctx, locationID, employeeIDs, start, end)` |
| `employeeRepository` | `GetByID(ctx, id)`, `GetByIDs(ctx, ids)` |
| `serviceRepository` | `GetByID(ctx, id)` |
| `employeeServiceRepository` | `GetByEmployeeAndService(...)`, `GetByLocationAndService(...)` |
| `locationRepository` | `GetByID(ctx, id)` |
| `scheduleRepository` | `GetLocationSlots(...)`, `GetEmployeesSlots(...)` |

> Репозитории нужны и для остальных booking-эндпоинтов (create, confirm, cancel), поэтому это общая задача.

---

## 5. Вопросы и решения

### Q1: Один эндпоинт vs два (как на staging)?

**Решение:** Один универсальный `POST /api/v1/bookings/slots` с диапазоном дат. Клиент сам выбирает — запросить месяц (для календаря дат) или один день (для тайм-слотов). Ответ всегда полный: с разбивкой по сотрудникам и конкретными слотами.

### Q2: Таймзоны

На main конвертация в Almaty (UTC+5) захардкожена в handler. На full-refactor нужно решить: хранить TZ в Location или оставить конвертацию на клиенте. Рекомендация — **все времена в UTC**, клиент конвертирует на своей стороне (TZ можно передавать в ответе Location).

### Q3: Совместимость со staging API

Новый API **не совместим** со staging (другие поля, UUID вместо point_code, один эндпоинт вместо двух). Это ожидаемо — full-refactor это полный рефакторинг.

---

## 6. Порядок задач

1. **UseCase `GetAvailableSlots`** — основная логика оркестрации
2. **Handler `get_slots.go`** — HTTP слой
3. **Роут + интерфейс** — подключение в handler.go
4. **Error mapping** — маппинг доменных ошибок
5. **Infra-репозитории** — PostgreSQL-реализации (общая задача для всего booking контекста)
6. **Тестирование** — unit-тесты для `generateSlots()` (алгоритм уже есть, тесты нужны)

Шаги 1-4 можно выполнить без инфра-репозиториев (код скомпилируется, но не запустится до реализации шага 5).
