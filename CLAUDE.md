

Ты — ведущий Backend-разработчик проекта Bagsy. 
Твоя цель — писать поддерживаемый, типобезопасный код на Go, следуя принципам Clean Architecture

# Руководство по разработке Bagsy

**Язык:** Go 1.23+ | **Архитектура:** Clean Architecture | **БД:** PostgreSQL (pgx/v5) | **Router:** chi/v5
**Core пакет:** `/home/rassulturtulov/Desktop/Programming/core`. Ссылка на репозиторий: `https://github.com/Rasikrr/core`

**Важные файлы:** `./notion_doc.pdf`, `./Makefile`, `./docs/swagger/`, `./.env.example`

---

## 1. АРХИТЕКТУРА

### Dependency Rule
```
HTTP Handlers → Services → Repositories → Domain
```

**ЗАПРЕЩЕНО:**
- Domain импортирует инфраструктуру (только stdlib)
- Services импортируют HTTP handlers
- Domain entities содержат db-теги

### DIP: Интерфейсы определяются где используются, ни где реализуются

**❌ НЕПРАВИЛЬНО:**
```go
// repositories/users/repository.go
type Repository interface { GetByPhone(...) }
type repository struct { ... }
```

**✅ ПРАВИЛЬНО:**
```go
// services/users/Service.go
type UserRepository interface { GetByPhone(...) }
type Service struct { repo UserRepository }

// handlers/auth/router.go
type UserCreator interface { CreateUser(...) }  // Минимальный интерфейс
type Controller struct { userService UserCreator }
```

**ISP:** Каждый слой определяет минимальный интерфейс. Разные handlers → разные интерфейсы для одного сервиса.

---

## 2. СТРУКТУРА СЛОЕВ

### Domain (`internal/domain/`)
```
domain/
├── entity/    # Сущности без db-тегов
├── enum/      # enumer генерируемые типы
├── errors/    # Доменные ошибки
├── query/     # DTO для фильтров
└── session/   # Контекст пользователя
```

**Правила:**
- Только простые Go структуры
- Enum через `enumer`
- Используй `cockroachdb/errors`
- Все ошибки через `domainErr.New*Error()`

**Доменные ошибки:**
```go
// errors/type.go
type ErrorType string // NOT_FOUND, INVALID_INPUT, VALIDATION, UNAUTHORIZED, FORBIDDEN, CONFLICT, INTERNAL

type DomainError struct {
    Type    ErrorType
    Message string
    Cause   error
    Details map[string]interface{}
}

func NewNotFoundError(msg string, cause error) *DomainError
// + NewInvalidInputError, NewValidationError, NewUnauthorizedError,
//   NewForbiddenError, NewConflictError, NewInternalError

// errors/entities.go
var (
    ErrUserNotFound = NewNotFoundError("user(s) not found", nil)
    ErrBagsyNotFound = NewNotFoundError("bagsy not found", nil)
)
```

**Enum пример:**
```go
//go:generate enumer -type=Role -json -trimprefix Role -transform=snake
type Role int8
const (
    RoleUser Role = iota
    RoleStaff
)
```

---

### Repository (`internal/repositories/`)
```
repositories/users/
├── repository.go   # БЕЗ интерфейса!
├── model.go        # DB модели с db-тегами
├── statements.go   # SQL константы
└── dto.go          # Опционально для JSONB
```

**Паттерн:**
```go
// repository.go
type Repository struct { db *postgres.Postgres }
func NewRepository(db *postgres.Postgres) *Repository

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    var m model
    err := pgxscan.Get(ctx, r.db, &m, getUserByPhoneSQL, phone)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainErr.ErrUserNotFound.WithError(err)
        }
        return nil, domainErr.NewInternalError("failed to get user from db", err)
    }
    out, err := m.convert()
    if err != nil {
	    return nil, domainErr.NewInternalError("failed to get user from db", err)
    }
}

// model.go - DB модель с тегами
type model struct {
    Phone string `db:"phone"`
    Role  string `db:"role"`  // enum хранится как string
}

func (m model) convert() (*entity.User, error) {
    role, err := enum.RoleString(m.Role)
    if err != nil {
        return nil, errors.Wrap(err, "invalid role")
    }
    return &entity.User{Phone: m.Phone, Role: role}, nil
}

// statements.go
const getUserByPhoneSQL = `SELECT ... FROM users WHERE phone = $1 AND deleted_at IS NULL`
```

**Динамические запросы (squirrel):**
```go
func buildQuery(filter query.UserFilter) (string, []any, error) {
    builder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
        Select("...").From("users").Where(sq.Eq{"deleted_at": nil})

    if filter.NetworkCode != nil {
        builder = builder.Where(sq.Eq{"network_code": *filter.NetworkCode})
    }
    if len(filter.Roles) > 0 {
        roleStrs := lo.Map(filter.Roles, func(r enum.Role, _ int) string { return r.String() })
        builder = builder.Where(sq.Eq{"role": roleStrs})
    }
    return builder.ToSql()
}
```

**Правила:**
- Repository НЕ определяет интерфейс
- Всегда `pgxscan.Get/Select`
- Инфраструктурные ошибки → доменные
- `pgx.ErrNoRows` → `ErrUserNotFound`
- Soft delete через `deleted_at`

---

### Service (`internal/services/`)
```
services/users/
├── Service.go  # БЕЗ интерфейса Service!
└── models.go   # Опционально
```

**Паттерн:**
```go
type Service struct {
    usersRepo UserRepository  // Интерфейс определен ЗДЕСЬ
    txManager database.TXManager
}

// Интерфейсы зависимостей
type UserRepository interface {
    GetByPhone(ctx context.Context, phone string) (*entity.User, error)
    GetByParams(ctx context.Context, filter query.UserFilter) ([]*entity.User, error)
    Create/Update/Delete...
}

func NewService(repo UserRepository, txm database.TXManager) *Service

func (s *Service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    return s.usersRepo.GetByPhone(ctx, phone) // Просто пробрасываем
}

func (s *Service) Create(ctx context.Context, user *entity.User) error {
    existing, err := s.usersRepo.GetByPhone(ctx, user.Phone)
    if err != nil {
        if !domainErr.IsNotFound(err) {
            return err
        }
    }

    if existing != nil && existing.Active {
        return domainErr.NewConflictError("user already exists", nil)
    }

    return s.usersRepo.Create(ctx, user)
}
```

**Транзакции:**
```go
func (s *Service) CreateWithRelations(ctx context.Context, user *entity.User) error {
    txOpts := database.TXOptions{IsolationLevel: enum.IsoLevelReadCommited}

    return s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
        if err := s.usersRepo.Create(txCtx, user); err != nil {
            return err
        }
        if err := s.pointsRepo.Create(txCtx, point); err != nil {
            return err
        }
        return nil
    })
}
```

**Правила:**
- Service НЕ определяет интерфейс для себя
- ТОЛЬКО доменные ошибки
- Обычно просто пробрасывает ошибки из repo
- Проверка типов: `domainErr.IsNotFound()`, `IsConflict()` и т.д.

---

### HTTP Layer (`internal/ports/http/`)
```
ports/http/
├── handlers/auth|users|points|bagsies/
├── middlewares/
└── server.go
```

**Паттерн handler:**
```go
// handlers/auth/router.go
type AuthService interface {
    Login(ctx, phone, password string) (access, refresh string, err error)
}

type UserCreator interface {  // Минимальный интерфейс!
    CreateUser(ctx context.Context, user *entity.User) error
}

type Controller struct {
    authService AuthService
    userService UserCreator
}

func New(authSvc AuthService, userSvc UserCreator, authMW middlewares.AuthMiddleware) *Controller

func (c *Controller) Init(router *chi.Mux) {
    router.Route("/api/v1/auth", func(r chi.Router) {
        r.Post("/login", c.login)
        r.Post("/register", c.authMiddleware.Handle(c.register))
    })
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    ctx := r.Context()

    if err := httputil.GetData(r, &req); err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    if err := req.validate(); err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    access, refresh, err := c.authService.Login(ctx, req.Phone, req.PasswordHash)
    if err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    httputil.SendData(ctx, w, loginResponse{access, refresh}, http.StatusOK)
}
```

**Models:**
```go
// models.go
//go:generate easyjson -all models.go
type loginRequest struct {
    Phone    string `json:"phone" validate:"required,min=10"`
    PasswordHash string `json:"password" validate:"required"`
}

func (r *loginRequest) validate() error {
    return GetValidator().Struct(r)
}
```

**Validator (singleton):**
```go
var (validate = validator.New(); validatorOnce sync.Once)

func GetValidator() *validator.Validate {
    validatorOnce.Do(func() {
        validate.RegisterValidation("valid_role_not_admin", validRoleNotAdminValidator)
    })
    return validate
}
```

**Swagger (обязательно!):**
**Всегда пиши Swagger описание ручек**
```go
// @Summary Авторизация пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Данные"
// @Success 200 {object} api.SuccessResponse{data=loginResponse}
// @Failure 400,401 {object} api.ErrorResponse
// @Router /api/v1/auth/login [post]
func (c *Controller) login(w, r) { ... }
```

**Правила:**
- Handler определяет минимальный интерфейс (ISP)
- Валидация на уровне HTTP
- Всегда `errors.HandleError(ctx, w, err)`
- Всегда пиши Swagger → `make swagger`

---

## 3. ОБРАБОТКА ОШИБОК

### Иерархия
```
pgx/http ошибки → domainErr.New*() → Service → errors.HandleError() → HTTP коды
```

**Маппинг в HTTP:**
```
TypeNotFound      → 404
TypeInvalidInput  → 400
TypeValidation    → 400
TypeUnauthorized  → 401
TypeForbidden     → 403
TypeConflict      → 409
TypeInternal      → 500
```

**По слоям:**

**Repository:** Оборачиваем инфраструктурные
```go
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domainErr.ErrUserNotFound.WithError(err)
}
return nil, domainErr.NewInternalError("failed to get user", err)
```

**Service:** Пробрасываем или проверяем тип
```go
user, err := s.repo.GetByPhone(ctx, phone)
if err != nil && !domainErr.IsNotFound(err) {
    return domainErr.NewInternalError("check failed", err)
}
```

**Handler:** Преобразуем в HTTP
```go
user, err := h.Service.GetByPhone(ctx, phone)
if err != nil {
    errors.HandleError(ctx, w, err)  // Автоматический маппинг
    return
}
```

**Проверка типов:**
```go
domainErr.IsNotFound(err)
domainErr.IsConflict(err)
domainErr.IsInternal(err)
// и т.д.
```

---

## 4. БАЗА ДАННЫХ

**Технологии:** pgx/v5, pgxscan, squirrel, goose

**Get:**
```go
var m model
err := pgxscan.Get(ctx, r.db, &m, sql, id)
if errors.Is(err, pgx.ErrNoRows) { return nil, domainErr.ErrNotFound.WithError(err) }
```

**Select:**
```go
var mm []model
err := pgxscan.Select(ctx, r.db, &mm, sql)
if errors.Is(err, pgx.ErrNoRows) { return []*entity.T{}, nil }  // Пустой список!
```

**Soft Delete:**
```go
UPDATE users SET deleted_at = NOW() WHERE id = ANY($1) AND deleted_at IS NULL
```

**Транзакции:**
```go
txOpts := database.TXOptions{IsolationLevel: enum.IsoLevelReadCommited}
s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error { ... })
```

**Уровни изоляции:** `IsoLevelReadCommited`, `IsoLevelRepeatableRead`, `IsoLevelSerializable`

---

## 5. DEPENDENCY INJECTION

```go
// app.go
func InitApp(ctx context.Context) (*App, error) {
    db := postgres.New(cfg)
    txManager := postgres.NewTXManager(db.Pool())

    // Repositories (конкретные типы)
    usersRepo := usersRepo.NewRepository(db)

    // Services (конкретные типы, удовлетворяют интерфейсам из Service)
    userSvc := usersService.NewService(usersRepo, txManager)
    authSvc := authService.NewService(usersRepo, redis, jwt)

    // Handlers (определяют интерфейсы, сервисы им удовлетворяют)
    authHandler := authHandlers.New(authSvc, userSvc, authMW)

    router := chi.NewRouter()
    authHandler.Init(router)

    return &App{router}, nil
}
```

---

## 6. СОГЛАШЕНИЯ

**Файлы:** `snake_case.go`, стандартные: `repository.go`, `Service.go`, `router.go`, `model.go`, `statements.go`

**Пакеты:** lowercase, множественное число: `users`, `points`, `bagsies`
❌ `users.UserService` → ✅ `users.Service`

**Переменные:** `PascalCase` (экспорт), `camelCase` (приватные)
**Интерфейсы:** БЕЗ "I" → ✅ `UserRepository` ❌ `IUserRepository`

**Импорты:**
```go
import (
    "context"
    "github.com/.../internal/domain/entity"
    domainErr "github.com/.../internal/domain/errors"
    "github.com/Rasikrr/core/database"
)
```

---

## 7. ГЕНЕРАЦИЯ

```go
//go:generate enumer -type=Role -json -trimprefix Role -transform=snake
//go:generate easyjson -all models.go
```

```bash
go generate ./...  # Все генераторы
make swagger       # Swagger docs
make lint          # golangci-lint
make migrate-up    # Применить миграции
```

---

## 8. ЧЕКЛИСТ КОММИТА

- [ ] `make lint` прошел
- [ ] Миграции применены
- [ ] Интерфейсы определены где используются (Service НЕ определяет для себя, Handler определяет минимальный)
- [ ] DB модели ≠ domain entities
- [ ] Инфраструктурные ошибки → доменные (repo), пробрасываются (Service), → HTTP (handler через `errors.HandleError`)
- [ ] Транзакции через `database.TXOptions`
- [ ] SQL в `statements.go`, soft delete
- [ ] Swagger добавлен → `make swagger`
- [ ] `go generate ./...` запущен

---

## 9. КОМАНДЫ

```bash
go generate ./...              # Генераторы
make test|lint|fmt|swagger     # Качество кода
make migrate-up|migrate-down   # Миграции
make run                       # Запуск
```

---

## 10. КЛЮЧЕВЫЕ ПРИНЦИПЫ

1. **DIP:** Интерфейсы там, где используются (Service НЕ определяет для себя, handlers определяют)
2. **ISP:** Каждый handler → минимальный интерфейс
3. **Единый тип ошибок:** ТОЛЬКО `internal/domain/errors`, никаких "сервисных"
4. **Потоки ошибок:** Repository оборачивает → Service пробрасывает → Handler → HTTP (через `errors.HandleError`)
5. **Разделение:** DB models ≠ domain entities
6. **Транзакции:** `database.TXOptions` + `enum.IsoLevel`
7. **Типобезопасность:** Enum везде вместо строк
8. **Soft delete:** Всегда через `deleted_at`

Документ живой — обновляй при появлении паттернов.