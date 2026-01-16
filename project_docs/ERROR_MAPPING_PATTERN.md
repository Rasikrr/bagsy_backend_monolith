# Паттерн маппинга ошибок для Services

## TL;DR - Быстрая шпаргалка

**Правило:** Service ВСЕГДА преобразует инфраструктурные ошибки → доменные

```go
// ❌ НЕПРАВИЛЬНО
func (s *Service) DoSomething() error {
    err := s.client.Call()
    if err != nil {
        return err  // Инфраструктурная ошибка напрямую!
    }
}

// ✅ ПРАВИЛЬНО
func (s *Service) DoSomething() error {
    err := s.client.Call()
    if err != nil {
        return s.mapClientError(err)  // Преобразуем!
    }
}
```

---

## Шаблон для любого Service

### 1. Импорты

```go
import (
    "github.com/.../internal/clients/<name>"  // Для доступа к ошибкам
    domainErr "github.com/.../internal/domain/errors"
    "github.com/cockroachdb/errors"  // Для errors.Is()
)
```

### 2. Функция маппинга

```go
func (s *Service) map<ClientName>Error(err error) error {
    // Validation errors → InvalidInput (400)
    if errors.Is(err, client.ErrEmptyField) {
        return domainErr.NewInvalidInputError("invalid data", err)
    }

    // Auth errors → Unauthorized (401)
    if errors.Is(err, client.ErrAuthFailed) {
        return domainErr.NewUnauthorizedError("auth failed")
    }

    // Not Found → NotFound (404)
    if errors.Is(err, client.ErrNotFound) {
        return domainErr.NewNotFoundError("resource not found", err)
    }

    // Business errors → Internal (500)
    if errors.Is(err, client.ErrBusinessLogic) {
        return domainErr.NewInternalError("business error", err)
    }

    // ОБЯЗАТЕЛЬНО: Fallback для неизвестных ошибок
    return domainErr.NewInternalError("client error", err)
}
```

### 3. Использование

```go
func (s *Service) CreateOrder(ctx context.Context, order *Order) error {
    err := s.paymentClient.Charge(ctx, order.Amount)
    if err != nil {
        return s.mapPaymentError(err)  // ✅ Мапим!
    }
    return nil
}
```

---

## Таблица маппинга

| Инфраструктурная ошибка | Доменная ошибка | HTTP | Когда использовать |
|-------------------------|-----------------|------|-------------------|
| `ErrEmpty*`, `ErrInvalid*` | `InvalidInputError` | 400 | Невалидные данные от клиента |
| `ErrAuthFailed`, `ErrUnauthorized` | `UnauthorizedError` | 401 | Проблемы с аутентификацией |
| `ErrForbidden` | `ForbiddenError` | 403 | Нет прав доступа |
| `ErrNotFound` | `NotFoundError` | 404 | Ресурс не найден |
| `ErrAlreadyExists` | `ConflictError` | 409 | Конфликт данных |
| Остальные | `InternalError` | 500 | Внутренние/сетевые/бизнес ошибки |

---

## Примеры для разных клиентов

### SMS/WhatsApp (Notifications)

```go
func (s *Service) mapSMSError(err error) error {
    if errors.Is(err, sms.ErrNoFunds) {
        return domainErr.NewInternalError("no funds to send SMS", err)
    }
    if errors.Is(err, sms.ErrTooManyRequests) {
        return domainErr.NewInternalError("rate limit", err)
    }
    return domainErr.NewInternalError("SMS failed", err)
}
```

### S3 (Storage)

```go
func (s *Service) mapS3Error(err error) error {
    if errors.Is(err, s3.ErrEmptyKey) {
        return domainErr.NewInvalidInputError("file key required", err)
    }
    if errors.Is(err, s3.ErrUploadFailed) {
        return domainErr.NewInternalError("upload failed", err)
    }
    return domainErr.NewInternalError("storage error", err)
}
```

### Database (Repository)

```go
func (s *Service) mapDBError(err error) error {
    if errors.Is(err, pgx.ErrNoRows) {
        return domainErr.NewNotFoundError("record not found", err)
    }
    if isPgUniqueViolation(err) {
        return domainErr.NewConflictError("already exists", err)
    }
    return domainErr.NewInternalError("database error", err)
}
```

### Payment Gateway

```go
func (s *Service) mapPaymentError(err error) error {
    if errors.Is(err, payment.ErrCardDeclined) {
        return domainErr.NewInvalidInputError("card declined", err)
    }
    if errors.Is(err, payment.ErrInsufficientFunds) {
        return domainErr.NewInvalidInputError("insufficient funds", err)
    }
    if errors.Is(err, payment.ErrInvalidCard) {
        return domainErr.NewInvalidInputError("invalid card", err)
    }
    return domainErr.NewInternalError("payment failed", err)
}
```

---

## Checklist

При добавлении нового инфраструктурного клиента:

- [ ] Создать `internal/clients/<name>/errors.go`
- [ ] Клиент возвращает СВОИ ошибки (НЕ доменные)
- [ ] Service создает `map<Name>Error()` функцию
- [ ] Использовать `errors.Is()` для проверки
- [ ] Каждая инфраструктурная ошибка → соответствующая доменная
- [ ] Обязательный fallback для неизвестных ошибок
- [ ] Сохранять оригинальную ошибку: `New*Error("msg", err)`

---

## Ключевые правила

### ✅ DO:

```go
// Мапить ВСЕ инфраструктурные ошибки
return s.mapClientError(err)

// Использовать errors.Is() (сохраняет stack trace)
if errors.Is(err, client.ErrFoo) { ... }

// Сохранять оригинальную ошибку
NewInternalError("msg", err)
//                      ^^^

// Иметь fallback
default:
    return domainErr.NewInternalError("client error", err)
```

### ❌ DON'T:

```go
// Возвращать инфраструктурные ошибки напрямую
return err  // ❌

// Использовать == (теряет stack trace)
if err == client.ErrFoo { ... }  // ❌

// Терять оригинальную ошибку
NewInternalError("msg", nil)  // ❌

// Забывать fallback
// Что если появится новая ошибка? Крэш!
```

---

## Почему это важно?

1. **Clean Architecture** - Domain изолирован от инфраструктуры
2. **Единообразие** - Handler всегда получает доменные ошибки
3. **Тестируемость** - Легко мокать и тестировать
4. **Гибкость** - Можно заменить провайдера без изменения domain
5. **Правильный HTTP** - Автоматический маппинг → HTTP коды

---

**Примеры:**
- ✅ `internal/services/notifications/service.go` - SMS/WhatsApp маппинг
- ✅ `internal/clients/USAGE_EXAMPLES.md` - Подробные примеры
- ✅ `internal/services/notifications/ERROR_MAPPING.md` - Детальная документация
