# Active & Permissions — Решения и открытые задачи

Дата: 2026-03-03

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
- `Active` фильтровать **в базе** (`WHERE active = true AND deleted_at IS NULL`), не в коде.
- `CanServeClients()` может зависеть только от `!IsDeleted() && Permissions.CanProvideServices`, без `Active` — деактивированный сотрудник и так не дойдёт до бронирования (middleware отсечёт).

### 2. Права на смену permissions (policy)

Текущее поведение в `policy.CanChangePermissions`:

| Кто меняет | Кому | Разрешено? |
|---|---|---|
| **Owner** | Любому, **включая себе** | Да |
| **Manager** | Staff своей локации | Да |
| **Manager** | Себе | Нет (Manager не Staff) |
| **Staff** | Кому угодно | Нет |

- Эндпоинт: `PATCH /employees/{id}/permissions`
- Owner меняет свои permissions через тот же эндпоинт, подставляя свой ID
- Self-modification для Owner — **намеренное** поведение (может убрать себя из бронирования)

## TODO

### ChangeRole должен сбрасывать permissions
- Сейчас `ChangeRole` не трогает permissions — при повышении staff→manager permissions остаются от старой роли
- **Вариант A (предпочтительный):** `ChangeRole` автоматически сбрасывает permissions на `DefaultPermissionsForRole(newRole)`
- **Вариант B:** `ChangePermissions` валидирует допустимые комбинации для роли

### Валидация permissions при ChangePermissions
- Сейчас принимает любые комбинации без привязки к роли
- Нужно решить: свободная настройка или ограничения по роли?
- Если жёсткие — добавить `Permissions.ValidateForRole(role)` в домен, вызывать из `SetPermissions`
- Если свободные — оставить как есть, но документировать

### Отвязать CanServeClients() от Active
- Сейчас: `IsActive() && Permissions.CanProvideServices` (где `IsActive = Active && !IsDeleted`)
- Предлагается: `!IsDeleted() && Permissions.CanProvideServices`
- Причина: middleware уже отсекает неактивных сотрудников, двойная проверка не нужна

### Добавить комментарий в policy
- В `CanChangePermissions` добавить комментарий что self-modification для Owner — намеренное поведение
