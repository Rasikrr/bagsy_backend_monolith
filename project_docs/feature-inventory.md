# Feature Inventory (Main Branch Analysis)

**Дата:** 19.02.2026
**Контекст:** Аудит перед рефакторингом и расширением функционала создания локаций.

---

## 1. Auth & Access Context (🔐 Готовность: 90%)

Полностью реализован цикл управления доступом для сотрудников (владельцев).

### Реализованные фичи:
- **Регистрация владельца (Owner Registration):**
    - `POST /register`: Создание временного запроса (Redis), отправка OTP.
    - `POST /register/verify`: Атомарное создание (Транзакция): Org + Owner + Trial Sub + Work History.
    - `POST /register/resend`: Повторная отправка кода (GreenAPI/GreenAPI Fallback).
- **Аутентификация (Authentication):**
    - `POST /login`: По номеру телефона и паролю (bcrypt).
    - `POST /refresh`: Ротация JWT (Access + Refresh в Redis).
    - `POST /logout`: Инвалидация Refresh-токена.
- **Сброс пароля (Password Reset):**
    - `POST /password/reset`: Запрос ссылки на WhatsApp/SMS.
    - `POST /password/reset/confirm`: Смена пароля через токен (Redis).
- **Контекст доступа (OrgContext):**
    - Middleware для сборки полного контекста (Employee + Org + Sub + Plan) за один SQL JOIN.
    - Механизм защиты от просроченных подписок на уровне бэкенда.

---

## 2. Billing Context (💳 Готовность: 70%)

Ядро биллинга готово, но автоматизация платежей и рекуррентные списания вынесены в будущие задачи.

### Реализованные фичи:
- **Тарифные планы (Plans):**
    - Поддержка трех уровней: `solo`, `point`, `network`.
    - Система ограничений (Capabilities): `max_locations`, `max_employees`, `max_services`.
- **Подписки (Subscriptions):**
    - Жизненный цикл статусов: `trial`, `active`, `past_due`, `suspended`, `canceled`.
    - Логика "CanOperate" (можно ли пользоваться сервисом).
    - Автоматический 2-месячный Trial при регистрации.

---

## 3. Location Context (📍 Готовность: 40%)

В процессе активной разработки.

### Реализованные фичи:
- **Управление категориями:**
    - Реалистичный сид из 10 категорий бизнеса (Beauty, Medical, Car Service и т.д.).
- **Создание локации:**
    - Валидация лимитов плана перед созданием.
    - Автоматическая генерация `Slug`.
    - Поддержка гео-координат (Latitude/Longitude) и адресов.
    - Транзакционный перенос владельца в первую созданную локацию.
    - Флаг `PromptOrgProfile` для UI (триггер настройки сети при открытии 2-й точки).

---

## 4. Identity Context (👤 Готовность: 50%)

Базовые сущности готовы, но CRUD-операции для персонала и клиентов находятся в зачатке.

### Реализованные фичи:
- **Сотрудники (Employees):**
    - Богатая модель `Employee` с поддержкой ролей (`owner`, `manager`, `staff`).
    - Система разрешений (Permissions): `can_provide_services`, `can_manage_schedule`.
    - Механизм трансфера между локациями.
- **Клиенты (Customers):**
    - Разделение на `Customer` (глобальный профиль по телефону) и `CustomerBase` (данные внутри конкретной организации).
    - Поддержка заметок (`CustomerNotes`) от сотрудников.

---

## 5. Infrastructure & Shared (🛠️ Готовность: 85%)

### Технический стек:
- **Transport:** HTTP (Chi), Swagger (swaggo), JSON (easyjson).
- **Database:** PostgreSQL (pgx, pgxscan), Redis (для сессий и OTP).
- **Integrations:** 
    - WhatsApp (GreenAPI) — основной канал.
    - SMS (SmsC/Mobizon) — fallback.
    - S3 (AWS SDK v2) — для медиа-ассетов.
- **DDD Utils:** 
    - `Phone`: Валидация и нормализация номеров (E.164).
    - `Money`: Работа с финансами (Decimal).
    - `Duration`: Работа с временными интервалами слотов.

---

## 6. Что НЕ реализовано (Ближайший Backlog)

1. **Catalog Context:** Создание услуг и категорий услуг (только доменные структуры).
2. **Schedule Context:** Работа с расписаниями точек и сотрудников.
3. **Booking Context:** Создание и управление записями (только доменные структуры).
4. **Staff Invite:** Флоу приглашения сотрудников (есть в планах, нет в `main`).
5. **Billing Workers:** Кронджоба для проверки просроченных подписок и автопродления.
