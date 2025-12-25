# Сравнение Entity и Repository Models

## 1. Network

### Entity (domain/entity/network.go)
```go
type Network struct {
	Code        string      // PK
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	UpdatedBy   string
}
```

### Model (repositories/networks/model.go)
```go
type model struct {
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`   // ⚠️ DIFF: *string vs string
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UpdatedBy   string     `db:"updated_by"`
}
```

**Различия:**
- `Description`: в entity это `string`, в model это `*string`

---

## 2. PointCategory

### Entity (domain/entity/point_category.go)
```go
type PointCategory struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string      // ⚠️ СУЩЕСТВУЕТ в entity
}
```

### Model (repositories/point_categories/model.go)
```go
type model struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`   // ⚠️ DIFF: *string vs string
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	// ⚠️ НЕТ UpdatedBy в model, но есть в entity
}
```

**Различия:**
- `Description`: в entity это `string`, в model это `*string`
- `UpdatedBy`: есть в entity (*string), НЕТ в model

---

## 3. ServiceCategory

### Entity (domain/entity/service_category.go)
```go
type ServiceCategory struct {
	ID          int
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   *string
}
```

### Model (repositories/service_categories/model.go)
```go
type model struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`   // ⚠️ DIFF: *string vs string
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	UpdatedBy   *string    `db:"updated_by"`
}
```

**Различия:**
- `Description`: в entity это `string`, в model это `*string`

---

## 4. ServiceSubcategory

### Entity (domain/entity/service_subcategory.go)
```go
type ServiceSubcategory struct {
	ID                int
	ServiceCategoryID int
	Name              string
	Description       string
	CreatedAt         time.Time
	UpdatedAt         *time.Time
	UpdatedBy         *string
}
```

### Model (repositories/service_subcategory/model.go)
```go
type model struct {
	ID                int        `db:"id"`
	Name              string     `db:"name"`
	Description       *string    `db:"description"`   // ⚠️ DIFF: *string vs string
	ServiceCategoryID int        `db:"service_category_id"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         *time.Time `db:"updated_at"`
	UpdatedBy         *string    `db:"updated_by"`
}
```

**Различия:**
- `Description`: в entity это `string`, в model это `*string`

---

## 5. Point

### Entity (domain/entity/point.go)
```go
type Point struct {
	Code        string
	Name        string
	Description string
	NetworkCode string
	CategoryID  int
	Address     Address         // ⚠️ вложенная структура
	City        string
	Active      bool
	Schedule    []Schedule      // ⚠️ массив структур
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	UpdatedBy   string
}

type Address struct {
	Coordinates Coordinates
	Street      string
	City        string
}

type Coordinates struct {
	Latitude  float64
	Longitude float64
}

type Schedule struct {
	WeekDay int
	Open    time.Time
	Close   time.Time
	AllDay  bool
	Comment string
}
```

### Model (repositories/points/model.go)
```go
type model struct {
	Code        string     `db:"code"`
	Name        string     `db:"name"`
	Description *string    `db:"description"`   // ⚠️ DIFF: *string vs string
	NetworkCode string     `db:"network_code"`
	CategoryID  int        `db:"category_id"`
	Address     []byte     `db:"address"`       // ⚠️ DIFF: []byte vs Address
	City        string     `db:"city"`
	Active      bool       `db:"active"`
	Schedule    []byte     `db:"schedule"`      // ⚠️ DIFF: []byte vs []Schedule
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UpdatedBy   string     `db:"updated_by"`
}
```

**Различия:**
- `Description`: в entity это `string`, в model это `*string`
- `Address`: в entity это `Address` (struct), в model это `[]byte` (JSONB)
- `Schedule`: в entity это `[]Schedule` (slice of structs), в model это `[]byte` (JSONB)
- Есть отдельные DTO для конвертации в `dto.go` (addressDTO, scheduleDTO)

---

## 6. Service

### Entity (domain/entity/service.go)
```go
type Service struct {
	ID              uuid.UUID
	PointCode       string
	CategoryID      int
	SubcategoryID   *int
	Name            string
	Description     *string
	DurationMinutes int
	Active          bool
	CreatedAt       time.Time
	UpdatedAt       *time.Time
	UpdatedBy       *string
}
```

### Model (repositories/services/model.go)
```go
type model struct {
	ID              uuid.UUID  `db:"id"`
	PointCode       string     `db:"point_code"`
	CategoryID      int        `db:"category_id"`
	SubcategoryID   *int       `db:"subcategory_id"`
	Name            string     `db:"name"`
	Description     *string    `db:"description"`
	DurationMinutes int        `db:"duration_minutes"`
	Active          bool       `db:"active"`
	CreatedAt       time.Time  `db:"created_at"`
	UpdatedAt       *time.Time `db:"updated_at"`
	UpdatedBy       *string    `db:"updated_by"`
}
```

**Различия:**
- ✅ ПОЛНОЕ СОВПАДЕНИЕ

---

## 7. MasterService

### Entity (domain/entity/master_service.go)
```go
type MasterService struct {
	ID          uuid.UUID
	MasterPhone string
	ServiceID   uuid.UUID
	Price       decimal.Decimal
	Active      bool
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	// ⚠️ НЕТ UpdatedBy
}
```

### Model (repositories/master_services/model.go)
```go
type model struct {
	ID          uuid.UUID       `db:"id"`
	MasterPhone string          `db:"master_phone"`
	ServiceID   uuid.UUID       `db:"service_id"`
	Price       decimal.Decimal `db:"price"`
	Active      bool            `db:"active"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   *time.Time      `db:"updated_at"`
	// ⚠️ НЕТ UpdatedBy
}
```

**Различия:**
- ✅ ПОЛНОЕ СОВПАДЕНИЕ
- ⚠️ ВНИМАНИЕ: В схеме БД НЕТ поля updated_by для master_services

---

## 8. Bagsy

### Entity (domain/entity/bagsy.go)
```go
type Bagsy struct {
	ID          uuid.UUID
	ServiceID   uuid.UUID
	PointCode   string
	ClientPhone string
	MasterPhone string
	Status      enum.BagsyStatus        // ⚠️ enum
	Price       decimal.Decimal
	StartAt     time.Time
	EndAt       time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	UpdatedBy   string
}
```

### Model (repositories/bagsies/model.go)
```go
type model struct {
	ID          uuid.UUID       `db:"id"`
	PointCode   string          `db:"point_code"`
	ClientPhone string          `db:"client_phone"`
	Status      string          `db:"status"`        // ⚠️ DIFF: string vs enum
	MasterPhone string          `db:"master_phone"`
	Price       decimal.Decimal `db:"price"`
	ServiceID   uuid.UUID       `db:"service_id"`
	StartAt     time.Time       `db:"start_at"`
	EndAt       time.Time       `db:"end_at"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   *time.Time      `db:"updated_at"`
	UpdatedBy   string          `db:"updated_by"`
}
```

**Различия:**
- `Status`: в entity это `enum.BagsyStatus`, в model это `string` (конвертация через enum.BagsyStatusString)

---

## 9. User

### Entity (domain/entity/user.go)
```go
type User struct {
	Phone       string
	Password    string
	Role        enum.Role        // ⚠️ enum
	Name        string
	Surname     string
	PointCode   *string
	NetworkCode *string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
	DeletedAt   *time.Time
	UpdatedBy   string
}
```

### Model (repositories/users/model.go)
```go
type model struct {
	Phone       string     `db:"phone"`
	Password    string     `db:"password"`
	Role        string     `db:"role"`        // ⚠️ DIFF: string vs enum
	Name        string     `db:"name"`
	Surname     string     `db:"surname"`
	PointCode   *string    `db:"point_code"`
	NetworkCode *string    `db:"network_code"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UpdatedBy   string     `db:"updated_by"`
}
```

**Различия:**
- `Role`: в entity это `enum.Role`, в model это `string` (конвертация через enum.RoleString)

---

## СВОДКА РАЗЛИЧИЙ

### 1. Description поля (string vs *string)
**Проблема**: В БД поле nullable, но в entity не pointer

Затронутые entity:
- Network: entity.Description (string) → model.Description (*string)
- PointCategory: entity.Description (string) → model.Description (*string)
- ServiceCategory: entity.Description (string) → model.Description (*string)
- ServiceSubcategory: entity.Description (string) → model.Description (*string)
- Point: entity.Description (string) → model.Description (*string)

**Рекомендация**: Изменить в entity на `*string` для consistency

---

### 2. Enum поля (enum vs string)
**Это ОК**: Правильный паттерн для работы с БД

Затронутые entity:
- User: entity.Role (enum.Role) → model.Role (string)
- Bagsy: entity.Status (enum.BagsyStatus) → model.Status (string)

---

### 3. JSONB поля (struct vs []byte)
**Это ОК**: Правильный паттерн для JSONB

Затронутые entity:
- Point: entity.Address (Address) → model.Address ([]byte)
- Point: entity.Schedule ([]Schedule) → model.Schedule ([]byte)

---

### 4. Отсутствующие поля

**PointCategory**:
- entity.UpdatedBy (*string) - ЕСТЬ
- model.UpdatedBy - НЕТ
- ⚠️ ПРОВЕРИТЬ: Есть ли updated_by в схеме БД для point_categories?

**Схема БД (migrations/20251011193901_schema.sql)**:
```sql
CREATE TABLE IF NOT EXISTS point_categories (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
    -- ⚠️ НЕТ updated_by!
);
```

**Рекомендация**:
- Убрать UpdatedBy из entity.PointCategory
- ИЛИ добавить updated_by в схему БД и model

---

## ИТОГО

### Критичные несоответствия:
1. ❌ **PointCategory.UpdatedBy** - есть в entity, нет в БД и model
2. ⚠️ **Description поля** - nullable в БД (*string в model), но не nullable в entity (string)

### Некритичные (паттерны конвертации):
3. ✅ **Enum поля** - правильная конвертация enum ↔ string
4. ✅ **JSONB поля** - правильная конвертация struct ↔ []byte

### Полные совпадения:
- ✅ Service
- ✅ MasterService (но нет updated_by в обоих - это норма)
