# Руководство по разработке для проекта Bagsy

Этот документ содержит правила, паттерны и соглашения для разработки в проекте [Bagsy](https://github.com/Rasikrr/bagsy_backend_monolith).

## Базовая информация

**Язык:** Go 1.23+
**Архитектура:** Clean Architecture
**Основной пакет:** [Core](https://github.com/Rasikrr/core) (локальная версия: `/home/rassulturtulov/Desktop/Programming/core`)
**База данных:** PostgreSQL (pgx/v5)
**HTTP Router:** chi/v5

**Важные файлы:**
- Документация проекта: `./notion_doc.pdf`
- Скрипты: `./Makefile`, `./scripts/`
- API документация: `./docs/swagger/`
- Переменные окружения: `./.env.example`
- Зависимости: `./go.mod`

---

## 1. АРХИТЕКТУРА: Clean Architecture

### Правило зависимостей (Dependency Rule)

Зависимости всегда направлены внутрь: от внешних слоев к внутренним.

```
HTTP Handlers → Services → Repositories → Domain
     ↓             ↓            ↓           ↓
Infrastructure  Application  Infrastructure  Core
```

**Запрещено:**
- Domain слой НЕ должен импортировать ничего кроме стандартной библиотеки и утилит
- Services НЕ должны импортировать HTTP handlers
- Services НЕ должны знать о существовании конкретных repository реализаций
- Domain entities НЕ должны содержать db-теги или другие инфраструктурные детали

### Принцип инверсии зависимостей (DIP)

**КРИТИЧЕСКИ ВАЖНО:** Интерфейсы определяются там, где используются, а НЕ там, где реализуются!

**❌ НЕПРАВИЛЬНО:**
```go
// В пакете реализации (repositories/users)
package users

type Repository interface {
    GetByPhone(ctx context.Context, phone string) (*entity.User, error)
    Create(ctx context.Context, user *entity.User) error
}

type repository struct { ... }
```

**✅ ПРАВИЛЬНО:**
```go
// В пакете использования (services/users)
package users

// Интерфейс определяется ЗДЕСЬ - где используется
type UserRepository interface {
    GetByPhone(ctx context.Context, phone string) (*entity.User, error)
    Create(ctx context.Context, user *entity.User) error
}

type service struct {
    usersRepo UserRepository // Зависимость от интерфейса
}

func NewService(usersRepo UserRepository) *service {
    return &service{usersRepo: usersRepo}
}
```

### Interface Segregation Principle (ISP)

Каждый слой определяет **минимальный** интерфейс, который ему нужен.

**Пример:**

```go
// ========== СЕРВИС (Полная реализация) ==========
// internal/services/users/service.go
package users

type UserService struct {
    repo UserRepository
}

func (s *UserService) CreateUser(email string) error { ... }
func (s *UserService) GetUserList() ([]string, error) { ... }
func (s *UserService) ChangePassword(id, pass string) error { ... }
func (s *UserService) DeleteUser(id string) error { ... }


// ========== HANDLER AUTH (Использует МИНИМУМ методов) ==========
// internal/ports/http/handlers/auth/router.go
package auth

// Интерфейс определяется ТУТ - только нужные методы!
type UserCreator interface {
    CreateUser(email string) error
}

type AuthHandler struct {
    userService UserCreator // Зависимость от узкого интерфейса
}

func NewAuthHandler(userService UserCreator) *AuthHandler {
    return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    _ = h.userService.CreateUser("test@example.com")
    w.WriteHeader(http.StatusCreated)
}

// ========== HANDLER ADMIN (Использует ДРУГИЕ методы) ==========
// internal/ports/http/handlers/admin/router.go
package admin

// Другой handler - другой интерфейс
type UserManager interface {
    GetUserList() ([]string, error)
    DeleteUser(id string) error
}

type AdminHandler struct {
    userService UserManager
}

func NewAdminHandler(userService UserManager) *AdminHandler {
    return &AdminHandler{userService: userService}
}

func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    users, _ := h.userService.GetUserList()
    // отправка ответа
}
```

**Преимущества такого подхода:**
1. Каждый слой зависит только от того, что ему действительно нужно
2. Легче тестировать (меньше методов для мока)
3. Явная документация того, что использует каждый компонент
4. Изоляция изменений - изменения в сервисе не влияют на handlers, если интерфейс не меняется

---

## 2. СТРУКТУРА СЛОЕВ

### Domain Layer (`internal/domain/`)

**Назначение:** Содержит бизнес-сущности, правила и доменные ошибки.

**Структура:**
```
internal/domain/
├── entity/             # Доменные сущности (User, Bagsy, Point и т.д.)
├── enum/               # Перечисления (Role, BagsyStatus)
├── errors/             # Доменные ошибки
├── query/              # DTO для сложных запросов (фильтры)
├── command/            # DTO для команд (создание, обновление)
└── session/            # Контекст сессии пользователя
```

**Правила:**
- Сущности - это простые Go структуры без тегов БД
- Используй `cockroachdb/errors` для работы с ошибками
- Все доменные ошибки определяются в `internal/domain/errors/`
- Все неизвестные ошибки оборачиваются в `domainErr.NewInternalError`
- Enum-типы генерируются через `enumer`

**Пример сущности:**
```go
// internal/domain/entity/user.go
package entity

import "time"
import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type User struct {
    Phone       string
    Password    *string
    Role        enum.Role
    Name        *string
    Surname     *string
    PointCode   *string
    NetworkCode *string
    Active      bool
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
}
```

**Доменные ошибки:**
```go
// internal/domain/errors/type.go
package errors

type ErrorType string

const (
    TypeNotFound     ErrorType = "NOT_FOUND"
    TypeInvalidInput ErrorType = "INVALID_INPUT"
    TypeValidation   ErrorType = "VALIDATION"
    TypeUnauthorized ErrorType = "UNAUTHORIZED"
    TypeForbidden    ErrorType = "FORBIDDEN"
    TypeConflict     ErrorType = "CONFLICT"
    TypeInternal     ErrorType = "INTERNAL"
)

type DomainError struct {
    Type    ErrorType
    Message string
    Cause   error
    Details map[string]interface{}
}

func (e *DomainError) Error() string { ... }
func (e *DomainError) WithError(err error) *DomainError { ... }

// Конструкторы
func NewNotFoundError(message string, cause error) *DomainError { ... }
func NewInvalidInputError(message string, cause error) *DomainError { ... }
func NewValidationError(message string, keyVals ...string) *DomainError { ... }
func NewUnauthorizedError(message string) *DomainError { ... }
func NewForbiddenError(message string) *DomainError { ... }
func NewConflictError(message string, cause error) *DomainError { ... }
func NewInternalError(message string, cause error) *DomainError { ... }
```

```go
// internal/domain/errors/entities.go
package errors

var (
    ErrUserNotFound    = NewNotFoundError("user(s) not found", nil)
    ErrBagsyNotFound   = NewNotFoundError("bagsy not found", nil)
    ErrNetworkNotFound = NewNotFoundError("network(s) not found", nil)
    ErrPointNotFound   = NewNotFoundError("point(s) not found", nil)
)
```

**Domain Query DTO (фильтры):**
```go
// internal/domain/query/user_filter.go
package query

import "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"

type UserFilter struct {
    NetworkCode *string
    PointCode   *string
    Roles       []enum.Role 
    Phones      []string
}
```

**Enums:**
```go
// internal/domain/enum/bagsy_status.go
package enum

//go:generate enumer -type=BagsyStatus -json -trimprefix BagsyStatus -transform=snake -output bagsy_status_enumer.go

type BagsyStatus uint8

const (
    BagsyStatusCreated BagsyStatus = iota
    BagsyStatusActive
    BagsyStatusCompleted
    BagsyStatusCanceled
)
```

---

### Repository Layer (`internal/repositories/`)

**Назначение:** Абстракция работы с БД и внешними хранилищами.

**Структура каждого репозитория:**
```
repositories/users/
├── repository.go      # Реализация Repository
├── model.go          # DB модели с db-тегами
├── statements.go     # SQL запросы как константы
└── dto.go           # DTO для JSONB и сложных типов (опционально)
```

**ВАЖНО:** Интерфейс Repository НЕ определяется в этом пакете! Он определяется в service, который его использует.

**Паттерн репозитория:**

```go
// ========== REPOSITORY (Реализация) ==========
// internal/repositories/users/repository.go
package users

import (
    "context"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
    domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
    "github.com/Rasikrr/core/database/postgres"
    "github.com/cockroachdb/errors"
    "github.com/georgysavva/scany/v2/pgxscan"
    "github.com/jackc/pgx/v5"
)

// Это просто структура с методами - БЕЗ интерфейса
type Repository struct {
    db *postgres.Postgres
}

func NewRepository(db *postgres.Postgres) *Repository {
    return &Repository{db: db}
}

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    var m model
    err := pgxscan.Get(ctx, r.db, &m, getUserByPhoneSQL, phone)
    if err != nil {
        // Оборачиваем инфраструктурную ошибку в доменную
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainErr.ErrUserNotFound.WithError(err)
        }
        return nil, domainErr.NewInternalError("failed to get user from db", err)
    }

    user, err := m.convert()
    if err != nil {
        return nil, domainErr.NewInternalError("failed to convert user model", err)
    }
    return user, nil
}

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
    m := convert(user)
    _, err := r.db.Exec(ctx, createUserSQL,
        m.Phone, m.Password, m.Role, m.Name, m.Surname,
        m.PointCode, m.NetworkCode)
    if err != nil {
        return domainErr.NewInternalError("failed to create user in db", err)
    }
    return nil
}

func (r *Repository) GetByParams(ctx context.Context, filter query.UserFilter) ([]*entity.User, error) {
    q, args, err := buildQuery(filter)
    if err != nil {
        return nil, domainErr.NewInternalError("failed to build query", err)
    }

    var mm models
    err = pgxscan.Select(ctx, r.db, &mm, q, args...)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return []*entity.User{}, nil // Пустой список, а не ошибка
        }
        return nil, domainErr.NewInternalError("failed to select users from db", err)
    }

    out, err := mm.convert()
	if err != nil {
		return nil, domainErr.NewInternalError("failed to convert user model", err)
	}
	return user, nil
}
```

```go
// internal/repositories/users/model.go
package users

import (
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
    "github.com/cockroachdb/errors"
)

// DB модель с тегами
type model struct {
    Phone       string     `db:"phone"`
    Password    *string    `db:"password"`
    Role        string     `db:"role"`
    Name        *string    `db:"name"`
    Surname     *string    `db:"surname"`
    PointCode   *string    `db:"point_code"`
    NetworkCode *string    `db:"network_code"`
    Active      bool       `db:"active"`
    CreatedAt   time.Time  `db:"created_at"`
    UpdatedAt   time.Time  `db:"updated_at"`
    DeletedAt   *time.Time `db:"deleted_at"`
}

type models []model

// Domain -> DB model
func convert(e *entity.User) model {
    return model{
        Phone:       e.Phone,
        Password:    e.Password,
        Role:        e.Role.String(),
        Name:        e.Name,
        Surname:     e.Surname,
        PointCode:   e.PointCode,
        NetworkCode: e.NetworkCode,
        Active:      e.Active,
        CreatedAt:   e.CreatedAt,
        UpdatedAt:   e.UpdatedAt,
        DeletedAt:   e.DeletedAt,
    }
}

// DB model -> Domain
func (m model) convert() (*entity.User, error) {
    role, err := enum.RoleString(m.Role)
    if err != nil {
        return nil, errors.Wrap(err, "invalid role in database")
    }

    return &entity.User{
        Phone:       m.Phone,
        Password:    m.Password,
        Role:        role,
        Name:        m.Name,
        Surname:     m.Surname,
        PointCode:   m.PointCode,
        NetworkCode: m.NetworkCode,
        Active:      m.Active,
        CreatedAt:   m.CreatedAt,
        UpdatedAt:   m.UpdatedAt,
        DeletedAt:   m.DeletedAt,
    }, nil
}

func (mm models) convert() ([]*entity.User, error) {
    users := make([]*entity.User, 0, len(mm))
    for _, m := range mm {
        user, err := m.convert()
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}
```

```go
// internal/repositories/users/statements.go
package users

const getUserByPhoneSQL = `
    SELECT phone, password, role, name, surname, point_code, network_code,
           active, created_at, updated_at, deleted_at
    FROM users
    WHERE phone = $1 AND deleted_at IS NULL
`

const createUserSQL = `
    INSERT INTO users (phone, password, role, name, surname, point_code, network_code)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
`
```

**Динамические запросы:**
```go
func buildQuery(filter query.UserFilter) (string, []any, error) {
    psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
    builder := psql.Select("phone", "password", "role", "name", "surname",
                           "point_code", "network_code", "active",
                           "created_at", "updated_at", "deleted_at").
        From("users").
        Where(sq.Eq{"deleted_at": nil})

    if filter.NetworkCode != nil {
        builder = builder.Where(sq.Eq{"network_code": *filter.NetworkCode})
    }

    if filter.PointCode != nil {
        builder = builder.Where(sq.Eq{"point_code": *filter.PointCode})
    }

    if len(filter.Roles) > 0 {
        // Конвертируем []enum.Role -> []string для SQL
        roleStrings := lo.Map(filter.Roles, func(role enum.Role, _ int) string {
            return role.String()
        })
        builder = builder.Where(sq.Eq{"role": roleStrings})
    }

    if len(filter.Phones) > 0 {
        builder = builder.Where(sq.Eq{"phone": filter.Phones})
    }

    return builder.ToSql()
}
```

**Правила:**
- Repository НЕ определяет свой интерфейс - это делает service
- Всегда используй `pgxscan` для сканирования результатов
- Для динамических запросов используй `squirrel`
- Используй soft delete через поле `deleted_at`
- Инфраструктурные ошибки ВСЕГДА оборачивай в доменные через `domainErr.NewInternalError`
- При `pgx.ErrNoRows` возвращай соответствующую доменную ошибку (ErrUserNotFound и т.д.)

---

### Service Layer (`internal/services/`)

**Назначение:** Бизнес-логика приложения, оркестрация операций.

**Структура каждого сервиса:**
```
services/users/
├── service.go     # ТОЛЬКО реализация (БЕЗ интерфейса Service)
└── models.go      # DTO для сложных параметров (опционально)
```

**ВАЖНО:** Сервисы работают ТОЛЬКО с доменными ошибками из `internal/domain/errors`. НЕТ отдельных "сервисных ошибок".

**КРИТИЧЕСКИ ВАЖНО:** Сервис НЕ определяет интерфейс `Service`! Интерфейсы определяются в handlers, которые используют сервис.

**Паттерн сервиса:**

```go
// ========== SERVICE (Только реализация) ==========
// internal/services/users/service.go
package users

import (
    "context"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/query"
    "github.com/Rasikrr/core/database"
)

// Это просто структура - БЕЗ интерфейса Service!
type Service struct {
    usersRepo UserRepository
    txManager database.TXManager
}

// Интерфейсы зависимостей определяются ЗДЕСЬ (в месте использования)
type UserRepository interface {
    GetByPhone(ctx context.Context, phone string) (*entity.User, error)
    GetByParams(ctx context.Context, filter query.UserFilter) ([]*entity.User, error)
    Create(ctx context.Context, user *entity.User) error
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, users ...*entity.User) error
}

func NewService(usersRepo UserRepository, txManager database.TXManager) *Service {
    return &Service{
        usersRepo: usersRepo,
        txManager: txManager,
    }
}

// Методы реализации - все публичные
func (s *Service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    user, err := s.usersRepo.GetByPhone(ctx, phone)
    if err != nil {
        return nil, err // Пробрасываем доменную ошибку
    }
    return user, nil
}

func (s *Service) Create(ctx context.Context, user *entity.User) error {
    if err := s.usersRepo.Create(ctx, user); err != nil {
        return err // Пробрасываем доменную ошибку
    }
    return nil
}

func (s *Service) GetByPointCode(ctx context.Context, pointCode string) ([]*entity.User, error) {
    users, err := s.usersRepo.GetByParams(ctx, query.UserFilter{
        PointCode: &pointCode,
    })
    if err != nil {
        return nil, err // Пробрасываем доменную ошибку
    }
    return users, nil
}

func (s *Service) GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.User, error) {
    users, err := s.usersRepo.GetByParams(ctx, query.UserFilter{
        NetworkCode: &networkCode,
    })
    if err != nil {
        return nil, err // Пробрасываем доменную ошибку
    }
    return users, nil
}

func (s *Service) GetByRoles(ctx context.Context, roles ...enum.Role) ([]*entity.User, error) {
    users, err := s.usersRepo.GetByParams(ctx, query.UserFilter{
        Roles: roles,
    })
    if err != nil {
        return nil, err // Пробрасываем доменную ошибку
    }
    return users, nil
}

func (s *Service) UpdatePassword(ctx context.Context, phone, newPassword string) error {
    // Бизнес-логика
    user, err := s.usersRepo.GetByPhone(ctx, phone)
    if err != nil {
        return err // Пробрасываем доменную ошибку
    }

    user.Password = &newPassword

    if err := s.usersRepo.Update(ctx, user); err != nil {
        return err // Пробрасываем доменную ошибку
    }
    return nil
}
```


**Работа с транзакциями (из core/database):**
```go
import (
    "github.com/Rasikrr/core/database"
    "github.com/Rasikrr/core/enum"
    domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

func (s *Service) CreateWithRelations(ctx context.Context, user *entity.User, point *entity.Point) error {
    txOpts := database.TXOptions{
        IsolationLevel: enum.IsoLevelReadCommited,
        ReadOnly:       false,
    }

    err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
        if err := s.usersRepo.Create(txCtx, user); err != nil {
            return err
        }

        if err := s.pointsRepo.Create(txCtx, point); err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return domainErr.NewInternalError("failed to create user with relations", err)
    }

    return nil
}
```

**Обработка специфичных ошибок в сервисе:**
```go
func (s *Service) Create(ctx context.Context, user *entity.User) error {
    // Проверяем существует ли пользователь
    existing, err := s.usersRepo.GetByPhone(ctx, user.Phone)
    if err != nil && !domainErr.IsNotFound(err) {
        // Если ошибка НЕ "not found" - что-то пошло не так
        return domainErr.NewInternalError("failed to check existing user", err)
    }

    // Если пользователь найден и активен - конфликт
    if existing != nil && existing.Active {
        return domainErr.NewConflictError("user already exists", nil)
    }

    // Создаем пользователя
    if err := s.usersRepo.Create(ctx, user); err != nil {
        return err // Пробрасываем доменную ошибку дальше
    }

    return nil
}
```

**Правила:**
- Сервис НЕ определяет интерфейс для себя
- Сервис определяет интерфейсы своих зависимостей (репозитории, кеши и т.д.)
- Используй **ТОЛЬКО** доменные ошибки из `internal/domain/errors`
- Для проверки типа ошибки используй `domainErr.IsNotFound()`, `domainErr.IsInternal()` и т.д.
- Доменные ошибки из репозиториев обычно пробрасываются дальше без изменений
- Для новых ошибок создавай предопределенные ошибки в `internal/domain/errors/`
- Для транзакций используй `database.TXOptions` с `enum.IsoLevel`

---

### HTTP Layer (`internal/ports/http/`)

**Назначение:** HTTP handlers, middleware, роутинг.

**Структура:**
```
ports/http/
├── handlers/
│   ├── auth/          # Аутентификация
│   ├── users/         # Управление пользователями
│   ├── points/        # Точки обслуживания
│   └── bagsies/       # Бронирования
├── middlewares/       # Middleware (auth, logging, etc.)
└── server.go          # Инициализация сервера
```

**КРИТИЧЕСКИ ВАЖНО:** Handler определяет интерфейс сервиса, который ему нужен (ISP).

**Паттерн handler:**

```go
// ========== HANDLER (Определяет свой интерфейс) ==========
// internal/ports/http/handlers/auth/router.go
package auth

import (
    "context"
    "net/http"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/session"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/middlewares"
    "github.com/go-chi/chi/v5"
	httputil "github.com/Rasikrr/core/http"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
)

// Интерфейс определяется ТУТ - только методы, нужные для auth
type AuthService interface {
    Login(ctx context.Context, phone, password string) (accessToken, refreshToken string, err error)
    CheckAccessToken(ctx context.Context, token string) (*session.Session, error)
    RefreshTokens(ctx context.Context, refreshToken string) (accessToken, refreshToken string, err error)
}

type UserCreator interface {
    CreateUser(ctx context.Context, user *entity.User) error
}

type Controller struct {
    authService    AuthService
    userService    UserCreator
    authMiddleware middlewares.AuthMiddleware
}

func New(authService AuthService, userService UserCreator, authMW middlewares.AuthMiddleware) *Controller {
    return &Controller{
        authService:    authService,
        userService:    userService,
        authMiddleware: authMW,
    }
}

func (c *Controller) Init(router *chi.Mux) {
    router.Route("/api/v1/auth", func(r chi.Router) {
        r.Post("/login", c.login)
        r.Post("/register", c.authMiddleware.Handle(c.register))
        r.Post("/refresh", c.refresh)
    })
}

func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
	ctx := r.Context()

    // Парсинг и валидация
    if err := httputil.GetData(r, &req); err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    if err := req.validate(); err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    // Вызов сервиса
    accessToken, refreshToken, err := c.authService.Login(r.Context(), req.Phone, req.Password)
    if err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    // Отправка ответа
	httputil.SendData(ctx, w, loginResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }), http.StatusOK)
}
```

```go
// internal/ports/http/handlers/admin/router.go
package admin

import (
    "context"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
	httputil "github.com/Rasikrr/core/http"
	"github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/errors"
)

// Другой handler - другой интерфейс!
type UserManager interface {
    GetByPointCode(ctx context.Context, pointCode string) ([]*entity.User, error)
    GetByNetworkCode(ctx context.Context, networkCode string) ([]*entity.User, error)
    GetByRoles(ctx context.Context, roles ...enum.Role) ([]*entity.User, error)
}

type AdminController struct {
    userService UserManager
}

func New(userService UserManager) *AdminController {
    return &AdminController{userService: userService}
}

func (c *AdminController) ListUsersByPoint(w http.ResponseWriter, r *http.Request) {
    pointCode := chi.URLParam(r, "pointCode")

    users, err := c.userService.GetByPointCode(r.Context(), pointCode)
    if err != nil {
        errors.HandleError(ctx, w, err)
        return
    }

    httputil.SendData(r.Context(), w, usersResponse{
		Users: users,
    }, http.StatusOK)
}
```

**Модели запросов/ответов:**
```go
// internal/ports/http/handlers/auth/models.go
package auth

//go:generate easyjson -all models.go

import "github.com/go-playground/validator/v10"

type loginRequest struct {
    Phone    string `json:"phone"    validate:"required,min=10,max=15"`
    Password string `json:"password" validate:"required"`
}

func (r *loginRequest) validate() error {
    return validator.New().Struct(r)
}

type loginResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}
```
** Валидатор **
** В каждом хендлере будут сови валидаторы и они будут инициализироваться один раз **
```go 
package auth

var (
	validate      = validator.New()
	validatorOnce sync.Once
)

func GetValidator() *validator.Validate {
	validatorOnce.Do(func() {
		validate.RegisterValidation("valid_role_not_admin", validRoleNotAdminValidator)
	})
	return validate
}

func validRoleNotAdminValidator(fl validator.FieldLevel) bool {
	field := fl.Field()

	// Validator автоматически разыменовывает указатели для кастомных валидаторов
	// Поэтому мы работаем со строкой напрямую
	if field.Kind() != reflect.String {
		log.Infof(context.Background(), "field is not a string, kind: %v", field.Kind())
		return false
	}

	value := field.String()

	// Проверяем, что это валидная роль из enum
	_, err := enum.RoleString(value)
	if err != nil {
		log.Infof(context.Background(), "role not found in enum: %s", value)
		return false // роль не найдена в enum
	}

	// Проверяем, что это не admin
	return value != enum.RoleAdmin.String()
}


```

**Swagger документация:**
** При написании хендлера всегда пиши swagger документацию и перегенерируй ее с помощью make swagger **
```go
// Login godoc
// @Summary Авторизация пользователя
// @Description Выполняет авторизацию пользователя по номеру телефона и паролю
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Данные для авторизации"
// @Success 200 {object} api.SuccessResponse{data=loginResponse}
// @Failure 400 {object} api.ErrorResponse
// @Failure 401 {object} api.ErrorResponse
// @Router /api/v1/auth/login [post]
func (c *Controller) login(w http.ResponseWriter, r *http.Request) {
    // implementation
}
```

**Правила:**
- Handler определяет минимальный интерфейс, который ему нужен (ISP)
- Разные handlers могут определять разные интерфейсы для одного и того же сервиса
- Валидация запросов на уровне HTTP
- Используй `easyjson` для оптимизации JSON

---

## 3. ОБРАБОТКА ОШИБОК

### Иерархия ошибок

```
Инфраструктурные ошибки (pgx, http, external APIs)
           ↓ (оборачиваем)
    Доменные ошибки (internal/domain/errors)
           ↓ (пробрасываем)
         Service
           ↓ (преобразуем через errors.HandleError (./internal/ports/http/errors))
    HTTP ответы (400, 401, 404, 500)
```

### Единственный тип ошибок: Доменные ошибки

**КРИТИЧЕСКИ ВАЖНО:** В проекте используются ТОЛЬКО доменные ошибки из `internal/domain/errors`. Нет отдельных "сервисных ошибок" или "HTTP ошибок".

**Доменные ошибки (`internal/domain/errors`)**

```go
import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

// Предопределенные ошибки
var (
    ErrUserNotFound    = NewNotFoundError("user(s) not found", nil)
    ErrBagsyNotFound   = NewNotFoundError("bagsy not found", nil)
    ErrNetworkNotFound = NewNotFoundError("network(s) not found", nil)
)

// Конструкторы для динамических ошибок
domainErr.NewNotFoundError("resource not found", cause)
domainErr.NewInvalidInputError("invalid data", cause)
domainErr.NewValidationError("validation failed", "field", "phone", "reason", "too short")
domainErr.NewUnauthorizedError("invalid credentials")
domainErr.NewForbiddenError("access denied")
domainErr.NewConflictError("user already exists", cause)
domainErr.NewInternalError("database error", cause)

// Оборачивание существующих ошибок
domainErr.ErrUserNotFound.WithError(pgxError)
```

### Правила обработки ошибок

**1. Используй `cockroachdb/errors` для технических операций:**
```go
import "github.com/cockroachdb/errors"

// Проверка типа
if errors.Is(err, pgx.ErrNoRows) { ... }

// Оборачивание (внутри пакета)
return errors.Wrap(err, "failed to convert model")
```

**2. Repository слой - оборачиваем инфраструктурные ошибки:**
```go
import (
    "github.com/jackc/pgx/v5"
    "github.com/cockroachdb/errors"
    domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

func (r *Repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    var m model
    err := pgxscan.Get(ctx, r.db, &m, getUserByPhoneSQL, phone)
    if err != nil {
        // ✅ ПРАВИЛЬНО - оборачиваем в доменную ошибку
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainErr.ErrUserNotFound.WithError(err)
        }
        return nil, domainErr.NewInternalError("failed to get user from db", err)
    }

    user, err := m.convert()
    if err != nil {
        return nil, domainErr.NewInternalError("failed to convert user model", err)
    }
    return user, nil
}

// ❌ НЕПРАВИЛЬНО - возвращаем сырую pgx ошибку
func (r *Repository) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    var m model
    err := pgxscan.Get(ctx, r.db, &m, getUserByPhoneSQL, phone)
    return &user, err // ← Плохо! Инфраструктурная ошибка утекает наружу
}
```

**3. Service слой - пробрасываем доменные ошибки:**
```go
import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

func (s *Service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    user, err := s.usersRepo.GetByPhone(ctx, phone)
    if err != nil {
        // ✅ ПРАВИЛЬНО - просто пробрасываем доменную ошибку
        return nil, err
    }
    return user, nil
}

// Если нужна дополнительная логика - проверяем тип ошибки
func (s *Service) Create(ctx context.Context, user *entity.User) error {
    existing, err := s.usersRepo.GetByPhone(ctx, user.Phone)
    if err != nil && !domainErr.IsNotFound(err) {
        // Если ошибка НЕ "not found" - что-то пошло не так
        return domainErr.NewInternalError("failed to check existing user", err)
    }

    // Если пользователь найден и активен - конфликт
    if existing != nil && existing.Active {
        return domainErr.NewConflictError("user already exists", nil)
    }

    if err := s.usersRepo.Create(ctx, user); err != nil {
        return err // Пробрасываем доменную ошибку
    }
    return nil
}

// ❌ НЕПРАВИЛЬНО - оборачивание в "сервисные ошибки" (их не существует!)
var errGetUser = NewError("failed to get user", 500) // ← Не делай так!

func (s *Service) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
    user, err := s.usersRepo.GetByPhone(ctx, phone)
    if err != nil {
        return nil, errGetUser.Wrap(err) // ← Плохо!
    }
    return user, nil
}
```

**4. Handler слой - преобразуем через errors.HandleError:**
```go
import "github.com/Rasikrr/core/api"

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    phone := chi.URLParam(r, "phone")

    user, err := h.userService.GetByPhone(r.Context(), phone)
    if err != nil {
        // ✅ ПРАВИЛЬНО - errors.HandleError автоматически преобразует доменную ошибку в HTTP код
        errors.HandleError(ctx, w, err)
        return
    }

    http.SendData(ctx, w, api.NewSuccessResponse(user), http.StatusOK)
}

// ❌ НЕПРАВИЛЬНО - вручную пишем JSON ответ
func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.userService.GetByPhone(r.Context(), phone)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError) // ← Плохо!
        json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
        return
    }
}
```

### Проверка типов доменных ошибок

```go
import domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

// Доступные проверки
domainErr.IsNotFound(err)       // TypeNotFound
domainErr.IsInvalidInput(err)   // TypeInvalidInput
domainErr.IsValidation(err)     // TypeValidation
domainErr.IsUnauthorized(err)   // TypeUnauthorized
domainErr.IsForbidden(err)      // TypeForbidden
domainErr.IsConflict(err)       // TypeConflict
domainErr.IsInternal(err)       // TypeInternal

// Пример использования
if domainErr.IsNotFound(err) {
    // Обработка "не найдено"
} else if domainErr.IsConflict(err) {
    // Обработка конфликта
} else {
    // Другие ошибки
}
```

### Маппинг доменных ошибок в HTTP коды

`errors.HandleError` автоматически преобразует доменные ошибки:

```
TypeNotFound      → 404 Not Found
TypeInvalidInput  → 400 Bad Request
TypeValidation    → 400 Bad Request
TypeUnauthorized  → 401 Unauthorized
TypeForbidden     → 403 Forbidden
TypeConflict      → 409 Conflict
TypeInternal      → 500 Internal Server Error
```

---

## 4. РАБОТА С БАЗОЙ ДАННЫХ

### Технологии

- **Драйвер:** `pgx/v5`
- **Сканирование:** `pgxscan`
- **Query builder:** `squirrel`
- **Миграции:** `goose`

### Паттерны работы с БД

**1. Простой запрос (Get):**
```go
func (r *Repository) GetByID(ctx context.Context, id string) (*entity.Network, error) {
    var m model
    err := pgxscan.Get(ctx, r.db, &m, getNetworkByIDSQL, id)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainErr.ErrNetworkNotFound.WithError(err)
        }
        return nil, domainErr.NewInternalError("failed to get network from db", err)
    }

    network, err := m.convert()
    if err != nil {
        return nil, domainErr.NewInternalError("failed to convert network model", err)
    }
    return network, nil
}
```

**2. Запрос списка (Select):**
```go
func (r *Repository) GetAll(ctx context.Context) ([]*entity.Network, error) {
    var mm []model
    err := pgxscan.Select(ctx, r.db, &mm, getAllNetworksSQL)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return []*entity.Network{}, nil // Пустой список, не ошибка
        }
        return nil, domainErr.NewInternalError("failed to get networks from db", err)
    }

    networks := make([]*entity.Network, 0, len(mm))
    for _, m := range mm {
        network, err := m.convert()
        if err != nil {
            return nil, domainErr.NewInternalError("failed to convert network model", err)
        }
        networks = append(networks, network)
    }
    return networks, nil
}
```

**3. Soft Delete:**
```go
func (r *Repository) Delete(ctx context.Context, ids ...string) error {
    _, err := r.db.Exec(ctx, softDeleteSQL, pq.Array(ids))
    if err != nil {
        return domainErr.NewInternalError("failed to soft delete records", err)
    }
    return nil
}

const softDeleteSQL = `
    UPDATE users
    SET deleted_at = NOW()
    WHERE id = ANY($1) AND deleted_at IS NULL
`
```

**4. Транзакции (через core/database):**
```go
import (
    "github.com/Rasikrr/core/database"
    "github.com/Rasikrr/core/enum"
    domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"
)

func (s *Service) CreateWithDependencies(ctx context.Context, data CreateData) error {
    txOpts := database.TXOptions{
        IsolationLevel: enum.IsoLevelReadCommited,
        ReadOnly:       false,
    }

    err := s.txManager.Transaction(ctx, txOpts, func(txCtx context.Context) error {
        if err := s.repo1.Create(txCtx, data.Entity1); err != nil {
            return err
        }

        if err := s.repo2.Create(txCtx, data.Entity2); err != nil {
            return err
        }

        return nil
    })

    if err != nil {
        return domainErr.NewInternalError("failed to create with dependencies", err)
    }

    return nil
}
```

**Уровни изоляции транзакций:**
```go
import "github.com/Rasikrr/core/enum"

// Доступные уровни
enum.IsoLevelReadCommited    // Read Committed (по умолчанию)
enum.IsoLevelRepeatableRead  // Repeatable Read
enum.IsoLevelSerializable    // Serializable
```

---

## 5. DEPENDENCY INJECTION

### Инициализация приложения

```go
// internal/app/app.go
package app

import (
    "context"
    "github.com/Rasikrr/core/database/postgres"

    // Repositories
    usersRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/users"
    pointsRepo "github.com/Rasikrr/bagsy_backend_monolith/internal/repositories/points"

    // Services
    usersService "github.com/Rasikrr/bagsy_backend_monolith/internal/services/users"
    authService "github.com/Rasikrr/bagsy_backend_monolith/internal/services/auth"

    // Handlers
    authHandlers "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/auth"
    adminHandlers "github.com/Rasikrr/bagsy_backend_monolith/internal/ports/http/handlers/admin"
)

func InitApp(ctx context.Context) (*App, error) {
    // 1. Инициализация инфраструктуры
    db := postgres.New(config.PostgresConfig)
    txManager := postgres.NewTXManager(db.Pool())

    // 2. Инициализация репозиториев (конкретные типы)
    usersRepository := usersRepo.NewRepository(db)
    pointsRepository := pointsRepo.NewRepository(db)

    // 3. Инициализация сервисов (конкретные типы)
    // Передаем конкретные реализации, они удовлетворяют интерфейсам из service
    userService := usersService.NewService(usersRepository, txManager)
    authSvc := authService.NewService(usersRepository, redisClient, jwtSecret, ...)

    // 4. Инициализация HTTP handlers
    // Handlers определяют свои интерфейсы, сервисы удовлетворяют им
    authHandler := authHandlers.New(authSvc, userService, authMiddleware)
    adminHandler := adminHandlers.New(userService)

    // 5. Регистрация роутов
    router := chi.NewRouter()
    authHandler.Init(router)
    adminHandler.Init(router)

    return &App{router: router}, nil
}
```

**Как это работает:**

1. **Repository** - конкретная структура `users.Repository`
2. **Service** - определяет интерфейс `UserRepository`, конкретный `users.Repository` удовлетворяет ему
3. **Handler** - определяет интерфейс `UserCreator` или `UserManager`, конкретный `users.Service` удовлетворяет им

Благодаря этому:
- Каждый слой зависит только от нужных методов (ISP)
- Легко тестировать (минимальные моки)
- Явная документация зависимостей

---

## 6. СОГЛАШЕНИЯ ПО ИМЕНОВАНИЮ

### Файлы
- **Snake_case:** `user_filter.go`, `bagsy_status.go`
- **Стандартные имена:**
  - `repository.go` - реализация репозитория
  - `service.go` - реализация сервиса
  - `router.go` - контроллер и роутинг
  - `model.go` - модели БД или DTO
  - `statements.go` - SQL запросы
  - `errors.go` - ошибки пакета
- **Тесты:** `*_test.go`
- **Генерированные:** `*_easyjson.go`, `*_enumer.go`

### Пакеты
- **Lowercase**, одно слово
- **Множественное число:** `users`, `points`, `bagsies`
- **Избегай статтера:**
  - ❌ `users.UserService`
  - ✅ `users.Service`

### Переменные и функции
- **Экспортируемые:** `PascalCase`
- **Неэкспортируемые:** `camelCase`
- **Интерфейсы:** БЕЗ префикса "I"
  - ✅ `UserRepository`, `UserCreator`
  - ❌ `IUserRepository`

### Импорты с алиасами
```go
import (
    "context"

    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/entity"
    "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/enum"
    domainErr "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/errors"

    "github.com/Rasikrr/core/api"
    "github.com/Rasikrr/core/database"
)
```

---

## 7. ГЕНЕРАЦИЯ КОДА

### Используемые генераторы

**1. Enumer (для enum-ов):**
```go
//go:generate enumer -type=Role -json -trimprefix Role -transform=snake -output role_enumer.go

type Role int8

const (
    RoleUser Role = iota
    RoleStaff
    RoleManager
)
```

**2. Easyjson (для JSON):**
```go
//go:generate easyjson -all models.go
```

**3. Swag (для Swagger):**
```bash
make swagger
```

**Запуск всех генераторов:**
```bash
go generate ./...
```
** Линтер (golangci-lit) **

```bash
make lint
```

** Применение миграций ** 
```bash 
make migrate-up
```
---

## 8. ЧЕКЛИСТ ПЕРЕД КОММИТОМ

- [ ] Линтеры прошли
- [ ] Миграции применены
- [ ] Интерфейсы определены там, где используются (ISP)
- [ ] Service НЕ определяет интерфейс для себя
- [ ] Handler определяет минимальный интерфейс сервиса
- [ ] Модели БД и доменные сущности разделены
- [ ] Инфраструктурные ошибки обернуты в доменные
- [ ] Сервисы работают ТОЛЬКО с доменными ошибками (без сервисных ошибок)
- [ ] Используется `cockroachdb/errors` для проверки типов
- [ ] Доменные ошибки преобразуются через `errors.HandleError` в handlers
- [ ] Транзакции используют `database.TXOptions` и `enum.IsoLevel`
- [ ] SQL запросы в `statements.go`
- [ ] Используется soft delete
- [ ] Swagger документация добавлена
- [ ] Генераторы запущены
- [ ] Тесты проходят
- [ ] Код отформатирован

---

## 9. ПОЛЕЗНЫЕ КОМАНДЫ

```bash
# Генерация кода
go generate ./...

# Запуск тестов
make test

# Обновление Swagger
make swagger

# Форматирование
make fmt

# Линтеры
make lint

# Миграции
make migrate-up
make migrate-down

# Запуск
make run
```

---

## 10. РЕЗЮМЕ: Ключевые принципы

1. **Интерфейсы там, где используются** - Service не определяет интерфейс для себя, это делают handlers
2. **ISP** - каждый handler определяет минимальный интерфейс
3. **Доменные ошибки** - ТОЛЬКО доменные ошибки из `internal/domain/errors` на всех слоях
4. **Repositories** - оборачивают инфраструктурные ошибки в доменные
5. **Services** - пробрасывают доменные ошибки без изменений
6. **Handlers** - преобразуют доменные ошибки в HTTP через `errors.HandleError`
7. **TXOptions** - используй `database.TXOptions` и `enum.IsoLevel`
8. **Enum везде** - вместо строк используй типобезопасные enum

Этот документ - живое руководство. Обновляй его при появлении новых паттернов.