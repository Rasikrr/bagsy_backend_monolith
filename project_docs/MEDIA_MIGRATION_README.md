# Media Tables Migration

## Созданные таблицы

### 1. `media` - Основная таблица медиа-файлов
Хранит информацию о всех файлах в S3.

**Поля:**
- `id` (UUID) - PRIMARY KEY
- `file_key` (TEXT) - Путь в S3, например: `2026/01/{uuid}.jpg`
- `bucket_name` (TEXT) - Имя S3 bucket (по умолчанию: 'bagsy-media')
- `original_filename` (TEXT) - Оригинальное имя файла
- `mime_type` (TEXT) - image/jpeg, image/png, image/webp
- `size_bytes` (BIGINT) - Размер файла (max 10MB)
- `width`, `height` (INTEGER) - Размеры изображения
- `status` (TEXT) - pending, active, processing, failed
- `uploaded_by_user_phone` (TEXT) - FK → users(phone)
- `created_at`, `updated_at`, `deleted_at` (TIMESTAMPTZ)

**Constraints:**
- ✅ Только разрешенные MIME types
- ✅ Размер файла: 1 byte - 10MB
- ✅ Dimensions: max 4096x4096

---

### 2. `user_media` - Аватары пользователей (1:1)

**Поля:**
- `user_phone` (TEXT) - PRIMARY KEY, FK → users(phone)
- `media_id` (UUID) - FK → media(id)
- `created_at`, `updated_at` (TIMESTAMPTZ)

**Особенность:**
- PRIMARY KEY на `user_phone` = только **один аватар** на пользователя

---

### 3. `staff_media` - Фото работ сотрудников (1:N)

**Поля:**
- `id` (UUID) - PRIMARY KEY
- `staff_phone` (TEXT) - FK → users(phone)
- `media_id` (UUID) - FK → media(id)
- `description` (TEXT) - Описание фото работы
- `display_order` (INTEGER) - Порядок в галерее
- `created_at`, `updated_at`, `deleted_at` (TIMESTAMPTZ)

**Constraints:**
- ✅ Уникальный `display_order` для каждого staff
- ✅ Soft delete

**Note:** staff = users с role='staff'

---

### 4. `point_media` - Фото локаций точек (1:N)

**Поля:**
- `id` (UUID) - PRIMARY KEY
- `point_code` (TEXT) - FK → points(code)
- `media_id` (UUID) - FK → media(id)
- `media_type` (TEXT) - exterior, interior, map, menu
- `is_primary` (BOOLEAN) - Главное фото для превью
- `display_order` (INTEGER) - Порядок отображения
- `created_at`, `updated_at`, `deleted_at` (TIMESTAMPTZ)

**Constraints:**
- ✅ Только **одно** primary фото на точку
- ✅ Уникальный `display_order` для каждой точки
- ✅ Soft delete

---

## Применение миграции

### 1. Убедитесь что POSTGRES_DSN настроен:
```bash
export POSTGRES_DSN="postgres://user:password@localhost:5432/bagsy?sslmode=disable"
```

Или в `.env`:
```
POSTGRES_DSN=postgres://user:password@localhost:5432/bagsy?sslmode=disable
```

### 2. Примените миграцию:
```bash
make migrate-up
```

### 3. Проверьте что таблицы созданы:
```bash
psql $POSTGRES_DSN -c "\dt media*"
psql $POSTGRES_DSN -c "\dt *_media"
```

---

## Откат миграции

**ВНИМАНИЕ:** Это удалит все данные в media таблицах!

```bash
# Откатить последнюю миграцию
make migrate-down-to VERSION=20260111120000

# Полный откат всех миграций (осторожно!)
make migrate-reset
```

---

## Архитектурные решения

### ✅ Почему гибридный подход (media + junction tables)?

1. **Referential Integrity** - FK constraints защищают от orphan records
2. **Бизнес-правила enforc'ятся БД:**
   - User: только 1 аватар → `PRIMARY KEY (user_phone)`
   - Point: только 1 primary фото → `UNIQUE INDEX WHERE is_primary = TRUE`
   - Уникальный порядок → `UNIQUE INDEX (entity_id, display_order)`
3. **Производительность** - простые JOIN'ы с индексами
4. **Каскадное удаление** - `ON DELETE CASCADE` автоматически чистит связи

### ✅ Почему не полиморфная таблица?

Полиморфная таблица (один `entity_media` для всех):
- ❌ Нет FK constraints
- ❌ Нельзя enforc'ить бизнес-правила
- ❌ Сложнее JOIN'ы
- ❌ Нет каскадного удаления

---

## Следующие шаги

1. ✅ Миграция создана
2. ⏳ Применить миграцию: `make migrate-up`
3. ⏳ Создать domain entity `Media`
4. ⏳ Создать repository `media`
5. ⏳ Добавить поле `AvatarURL` в `entity.User`
6. ⏳ Обновить `users` repository (LEFT JOIN media)
7. ⏳ Реализовать upload flow в handlers

---

## Структура S3

```
bucket: bagsy-media/
├── 2026/
│   ├── 01/
│   │   ├── 550e8400-e29b-41d4-a716-446655440000.jpg
│   │   └── 660e8400-e29b-41d4-a716-446655440001.png
│   └── 02/
│       └── ...
└── 2027/
    └── ...
```

**Формат file_key:** `YYYY/MM/{uuid}.{ext}`

---

**Дата создания:** 2026-01-11
**Файл миграции:** `20260111120000_create_media_table.sql`
