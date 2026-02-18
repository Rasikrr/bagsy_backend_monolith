# Password Change Flow (main branch)

Двухшаговый процесс сброса/смены пароля через one-time токен, отправляемый пользователю по WhatsApp (с fallback на SMS).

## Эндпоинты

| Шаг | Метод | URL | Rate Limiter | Auth |
|-----|-------|-----|:------------:|:----:|
| 1   | POST  | `/api/v1/auth/password/change` | 3 req/min | Нет |
| 2   | POST  | `/api/v1/auth/password/change/confirm` | — | Нет |

---

## Шаг 1 — Запрос на сброс пароля

**`POST /api/v1/auth/password/change`**

### Request Body

```json
{
  "phone": "77001234567"
}
```

| Поле  | Тип    | Валидация              |
|-------|--------|------------------------|
| phone | string | required, min=10, max=15 |

### Что происходит

```
Handler (change_password.go)
  │
  ├─ request.GetAndValidateData → парсинг + валидация changePasswordRequest
  │
  └─ authService.SendPasswordChangeLink(ctx, phone)
       │
       ├─ usersService.ExistsByPhone(ctx, phone)
       │   └─ usersRepo.ExistsByPhone → SELECT EXISTS ... FROM users WHERE phone = $1
       │   └─ Если пользователь не найден → return user.ErrUserNotFound (404)
       │
       ├─ codegen.GenerateAuthToken()
       │   └─ Генерирует случайный 10-символьный токен (a-z, 0-9)
       │
       ├─ tokensCache.SaveInviteToken(ctx, token, payload, registrationTTL)
       │   └─ Сохраняет в Redis:
       │        key:   token (10 chars)
       │        value: InviteTokenInfo{Phone, Purpose: TokenPurposePasswordChange}
       │        TTL:   registrationTTL (настраиваемый, обычно 24ч)
       │
       └─ notificationService.SendPasswordChangeLink(ctx, phone, token)
            └─ Формирует ссылку: "{registrationConfirmURL}/{token}"
            └─ Сообщение: "Для смены пароля в Bagsy следуйте по данной ссылке: {link}"
            └─ Отправка: WhatsApp → fallback SMS (если WhatsApp fail)
```

### Успешный ответ (200)

```json
{
  "message": "link sent"
}
```

### Ошибки

| HTTP | Когда |
|------|-------|
| 400  | Невалидный формат phone |
| 404  | Пользователь с таким phone не найден (`user.ErrUserNotFound`) |
| 500  | Ошибка Redis, ошибка отправки уведомления (оба канала упали) |

---

## Шаг 2 — Подтверждение смены пароля

**`POST /api/v1/auth/password/change/confirm`**

### Request Body

```json
{
  "token": "a1b2c3d4e5",
  "password": "newSecurePassword123"
}
```

| Поле     | Тип    | Валидация |
|----------|--------|-----------|
| token    | string | required  |
| password | string | required  |

### Что происходит

```
Handler (change_password_confirm.go)
  │
  ├─ request.GetAndValidateData → парсинг + валидация passwordChangeConfirmRequest
  │
  ├─ req.toDomain() → auth.ChangePasswordConfirmCommand{Token, Password}
  │
  └─ authService.ChangePassword(ctx, cmd)
       │
       ├─ tokensCache.GetInviteToken(ctx, token)
       │   └─ Ищет токен в Redis
       │   └─ Если не найден/истёк → return UnauthorizedError("invalid or expired token") (401)
       │   └─ Возвращает InviteTokenInfo{Phone, Purpose}
       │
       ├─ usersService.UpdatePasswordByPhone(ctx, phone, password)
       │   │
       │   ├─ usersRepo.GetByPhone(ctx, phone)
       │   │   └─ Загружает пользователя из PostgreSQL
       │   │
       │   ├─ hash.Password(rawPassword)
       │   │   └─ Хэширует пароль через bcrypt (pkg/hash)
       │   │
       │   └─ usersRepo.Update(ctx, user)
       │       └─ UPDATE users SET password_hash = $1 ... WHERE phone = $2
       │
       └─ tokensCache.DeleteInviteToken(ctx, token)
            └─ Удаляет использованный токен из Redis (one-time use)
            └─ Ошибка игнорируется — токен в любом случае истечёт по TTL
```

### Успешный ответ (200)

```json
{
  "message": "password changed"
}
```

### Ошибки

| HTTP | Когда |
|------|-------|
| 400  | Невалидные поля (пустой token или password) |
| 401  | Токен не найден, истёк или уже использован |
| 500  | Ошибка PostgreSQL, ошибка хэширования пароля |

---

## Хранилища данных

| Хранилище  | Что хранится | Ключ | TTL |
|------------|-------------|------|-----|
| Redis (tokensCache) | `InviteTokenInfo{Phone, Purpose: password_change}` | 10-char token | `registrationTTL` |
| PostgreSQL (users)  | Обновлённый `password_hash` | phone | — |

---

## Безопасность

| Аспект | Реализация |
|--------|------------|
| Токен | 10 символов, случайный (a-z, 0-9), one-time use |
| TTL токена | Ограничен `registrationTTL` в Redis |
| Rate limiting | Шаг 1 защищён rate limiter: **3 запроса/минуту** (по IP) |
| Хэширование пароля | bcrypt (`pkg/hash`) |
| Каналы доставки | WhatsApp → SMS fallback |
| Сессии | **Не инвалидируются** (swagger-док упоминает инвалидацию, но в коде refresh-токены не удаляются при смене пароля) |

---

## Схема потока (sequence)

```
Client                  Handler              AuthService         UsersService      Redis          Notification
  │                       │                      │                    │               │                │
  │── POST /change ──────>│                      │                    │               │                │
  │                       │── SendPasswordChangeLink ──>              │               │                │
  │                       │                      │── ExistsByPhone ──>│               │                │
  │                       │                      │<── true ───────────│               │                │
  │                       │                      │── GenerateAuthToken (10 chars) ──> │                │
  │                       │                      │── SaveInviteToken ────────────────>│                │
  │                       │                      │── SendPasswordChangeLink ─────────────────────────> │
  │                       │                      │                    │               │    WhatsApp/SMS │
  │<── 200 "link sent" ──│                      │                    │               │                │
  │                       │                      │                    │               │                │
  │── POST /change/confirm ─>│                   │                    │               │                │
  │   {token, password}   │── ChangePassword ───>│                   │               │                │
  │                       │                      │── GetInviteToken ─────────────────>│                │
  │                       │                      │<── InviteTokenInfo{phone} ─────────│                │
  │                       │                      │── UpdatePasswordByPhone ─>│        │                │
  │                       │                      │                    │ GetByPhone    │                │
  │                       │                      │                    │ hash.Password │                │
  │                       │                      │                    │ Update (PG)   │                │
  │                       │                      │<── ok ────────────│               │                │
  │                       │                      │── DeleteInviteToken ──────────────>│                │
  │<── 200 "password changed" ──│                │                    │               │                │
```

---

## Задействованные файлы (main branch)

| Слой | Файл | Роль |
|------|------|------|
| Handler | `internal/ports/http/handlers/auth/change_password.go` | Шаг 1 — HTTP handler |
| Handler | `internal/ports/http/handlers/auth/change_password_confirm.go` | Шаг 2 — HTTP handler |
| Handler | `internal/ports/http/handlers/auth/models.go` | DTO: `changePasswordRequest`, `passwordChangeConfirmRequest` |
| Handler | `internal/ports/http/handlers/auth/router.go` | Роутинг + rate limiter |
| Service | `internal/services/auth/service.go` | `SendPasswordChangeLink`, `ChangePassword` |
| Service | `internal/services/auth/types.go` | `InviteTokenInfo`, `TokenPurpose` |
| Service | `internal/services/users/service.go` | `UpdatePasswordByPhone`, `ExistsByPhone` |
| Service | `internal/services/notifications/service.go` | `SendPasswordChangeLink` (WhatsApp/SMS) |
| Domain  | `internal/domain/auth/command.go` | `ChangePasswordConfirmCommand` |
| Domain  | `internal/domain/auth/errors.go` | Sentinel errors |
| Domain  | `internal/domain/user/errors.go` | `ErrUserNotFound` |
| Util    | `internal/util/codegen/code.go` | `GenerateAuthToken` (10 chars) |
| Pkg     | `pkg/hash/hash.go` | bcrypt hashing |
| Errors  | `internal/ports/http/errors/` | `HandleError`, mapping domain → HTTP |

---

## Замечания

1. **Сессии не инвалидируются.** В swagger-документации к `changePasswordConfirm` сказано: *"все активные сессии (refresh токены) инвалидируются"*, однако в коде `ChangePassword` вызывается только `UpdatePasswordByPhone` и `DeleteInviteToken` — refresh-токены **не удаляются**. Старые access-токены продолжат работать до истечения TTL.
2. **Purpose не проверяется.** `ChangePassword` получает `InviteTokenInfo` из Redis, но не проверяет `Purpose == TokenPurposePasswordChange`. Теоретически, токен с `Purpose = Register` мог бы быть использован для смены пароля.
3. **Нет проверки активности пользователя.** В `SendPasswordChangeLink` проверяется только `ExistsByPhone`, но не `user.Active`. Неактивный пользователь может запросить смену пароля.
