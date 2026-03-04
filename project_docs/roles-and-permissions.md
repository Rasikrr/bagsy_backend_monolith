# Roles & Permissions — Полный отчёт

Дата: 2026-03-04

## 1. Роли

| Роль | Константа | Описание |
|------|-----------|----------|
| **Owner** | `RoleOwner = "owner"` | Владелец организации. Полный доступ. |
| **Manager** | `RoleManager = "manager"` | Менеджер локации. Управляет staff своей точки. |
| **Staff** | `RoleStaff = "staff"` | Мастер/сотрудник. Принимает клиентов. |

## 2. Permissions (настраиваемые)

| Поле | Описание |
|------|----------|
| `CanProvideServices` | Отображается при бронировании, принимает клиентов |
| `CanManageLocationSchedule` | Может управлять расписанием локации |

### Дефолтные permissions по ролям

| Роль | CanProvideServices | CanManageLocationSchedule |
|------|--------------------|---------------------------|
| Owner | `true` | `true` |
| Manager | `false` | `true` |
| Staff | `true` | `false` |

> При смене роли через `ChangeRole()` permissions автоматически сбрасываются на дефолтные для новой роли.

## 3. Три уровня отключения сотрудника

| Механизм | Семантика | Что блокирует |
|----------|-----------|---------------|
| `DeletedAt` | Уволен навсегда (soft delete) | Всё. Необратимо. |
| `Active = false` | Временно заблокирован | Логин, доступ к панели, все API. |
| `CanProvideServices = false` | Не принимает клиентов | Только видимость при бронировании. |

- `Active` проверяется при логине, refresh tokens и в middleware (SQL `WHERE active = true`).
- `CanServeClients()` = `!IsDeleted() && CanProvideServices` (не зависит от `Active`).

## 4. Матрица прав: кто что может

| Действие | Owner | Manager | Staff |
|----------|-------|---------|-------|
| **Приглашение сотрудника** | Любая роль кроме Owner | Только Staff | — |
| **Список сотрудников** | Вся организация | Только своя локация | — |
| **Activate / Deactivate** | Любой (кроме себя) | Staff своей локации | — |
| **Перевод (Transfer)** | Любой (кроме себя) | — | — |
| **Смена роли** | Любой (кроме себя) | — | — |
| **Смена permissions** | Любой, **включая себя** | Staff своей локации | — |
| **Создание локации** | Да (в лимитах плана) | — | — |
| **Отмена записи (Appointment)** | Любая в организации | Своей локации | Только свои |

## 5. Детали по каждой policy

### CanInviteEmployee

- Owner → может приглашать Manager и Staff (не Owner).
- Manager → может приглашать только Staff.
- Staff → не может приглашать.
- Проверяет лимит `MaxEmployees` из плана подписки.

### CanListEmployees

- Owner → видит всех сотрудников организации.
- Manager → filter автоматически устанавливается на `LocationID` менеджера.
- Staff → запрещено.

### CanManageEmployee (activate/deactivate)

- Owner → любого, кроме себя.
- Manager → только Staff своей локации (проверяет `targetEmp.LocationID == orgCtx.Employee.LocationID`).
- Staff → запрещено.
- Нельзя менять себя (`ErrCannotModifySelf`).

### CanTransferEmployee

- **Только Owner.** Менеджер не может переводить сотрудников.
- Нельзя переводить себя.
- Целевая локация должна принадлежать той же организации и быть активной.
- При трансфере создаётся запись в `WorkHistory` с `ChangeTypeTransfer`.

### CanChangeRole

- **Только Owner.**
- Нельзя менять свою роль.
- Нельзя назначить роль Owner (`ErrCannotSetOwnerRole`).
- При смене роли permissions сбрасываются на дефолтные.

### CanChangePermissions

- Owner → может менять permissions **любому, включая себе** (намеренно: может убрать себя из бронирования).
- Manager → только Staff своей локации.
- Staff → запрещено.

### CanCreateLocation

- **Только Owner.**
- Проверяет лимит `MaxLocations` из плана подписки.

### CanCancelAppointment

- Owner → любую запись в организации.
- Manager → только записи своей локации (`appt.BelongsToLocation`).
- Staff → только свои записи (`appt.BelongsToEmployee`).

## 6. Общие правила

1. **Subscription gate** — каждая policy проверяет `orgCtx.Subscription.Status.CanOperate()`. Если подписка приостановлена → `ErrSubscriptionSuspended`.
2. **Multi-tenancy** — все проверки гарантируют что `orgCtx.Organization.ID == targetEmp.OrganizationID`.
3. **Self-modification** — запрещена везде, кроме `CanChangePermissions` для Owner.
4. **Scope менеджера** — Manager видит и управляет только Staff **своей** локации.
