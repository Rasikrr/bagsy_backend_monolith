# Active & Permissions — Решения

Дата: 2026-03-03, обновлено 2026-03-04

## Контекст

Анализ полей `Employee.Active`, `Permissions` и их взаимодействия с бронированием (`CanServeClients`) и policy (`CanChangePermissions`).

## Решения

### 1. Active — оставляем

`Active` нужен. Это системный флаг доступа (логин, панель).

Три уровня отключения сотрудника:

| Механизм | Семантика | Что блокирует |
|---|---|---|
| `DeletedAt` | Уволен навсегда (soft delete) | Всё. Необратимо |
| `Active = false` | Временно заблокирован | Логин + вся видимость в системе |
| `CanProvideServices = false` | Не принимает клиентов | Только бронирование |

- Сценарий: конфликт/разбирательство — нужно забрать доступ, но не удалять сотрудника.
- `Active` фильтруется **в базе** (`WHERE active = true AND deleted_at IS NULL`) в `getOrgContext`.
- `Active` проверяется в `LoginEmployee()` и `RefreshTokens()`.
- `CanServeClients()` = `Active && !IsDeleted() && CanProvideServices` — Active нужен потому что booking публичный (без auth middleware).

### 2. Права на смену permissions (policy)

| Кто меняет | Кому | Разрешено? |
|---|---|---|
| **Owner** | Любому, **включая себе** | Да |
| **Manager** | Staff своей локации | Да |
| **Manager** | Себе | Нет (Manager не Staff) |
| **Staff** | Кому угодно | Нет |

- Self-modification для Owner — **намеренное** поведение (может убрать себя из бронирования)
- Задокументировано комментарием в `CanChangePermissions`

## Выполнено (2026-03-04)

- [x] **ChangeRole сбрасывает permissions** на `DefaultPermissionsForRole(newRole)` (Вариант A)
- [x] **Валидация permissions** — решено оставить свободную настройку без ограничений по роли
- [x] **CanServeClients()** — `Active && !IsDeleted() && CanProvideServices` (Active оставлен для публичных эндпоинтов)
- [x] **Комментарий в policy** — добавлен в `CanChangePermissions`
- [x] **Active проверяется при логине и refresh** — `IsActive()` в `LoginEmployee()` и `RefreshTokens()`
- [x] **OrgContext middleware** — SQL фильтрует `active = true`, проверяет `Organization.Active`
- [x] **Booking проверяет Location и Service** — `loc.CanOperate()`, `svc.IsActive()`
- [x] **Transfer — только Owner** — выделен в отдельную policy `CanTransferEmployee`
- [x] **Transfer пишет WorkHistory** — атомарно в транзакции
- [x] **Transfer валидирует локацию** — принадлежит организации + `CanOperate()`
