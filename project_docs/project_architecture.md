Вот выжимка архитектурных стандартов и подходов вашего проекта, основанная на предоставленной структуре папок. Сохраните этот текст как `ARCHITECTURE.md` или `README.md` в корне репозитория.

---

# Architecture & Design Guidelines

Этот проект построен на принципах **Clean Architecture** и **Domain-Driven Design (DDD)** в рамках **Модульного Монолита (Modular Monolith)**.

Наша цель — высокая связность внутри модулей (High Cohesion) и низкая связанность между ними (Low Coupling).

## 1. Глобальная Архитектура

Мы строго следуем **Dependency Rule** (Правилу Зависимостей): зависимости направлены только внутрь.

* **Core (Domain):** Ни от чего не зависит. Чистая бизнес-логика.
* **Application (Service):** Зависит от Domain. Оркестрирует бизнес-процессы.
* **Infrastructure (Repository, Gateway):** Зависит от Domain и Application. Реализует интерфейсы.
* **Presentation (Handler, API):** Зависит от Application.

### Поток данных

`HTTP Request` -> `Handler` -> `Service` -> `Domain Entity` + `Repository` -> `Database`

---

## 2. Структура Слоев (`internal/`)

### 🟢 1. Domain Layer (`internal/domain/`)

**Сердце системы.** Содержит бизнес-логику, независимую от фреймворков и БД.

* **Структура:** Разбита по бизнес-контекстам (`auth`, `identity`, `billing`, `booking`).
* **Содержимое:**
* **Aggregates & Entities:** (`employee.go`, `appointment.go`) — богатые модели с методами изменения состояния.
* **Value Objects:** (`money`, `phone`) — неизменяемые объекты.
* **Domain Errors:** (`errors.go`) — ошибки бизнес-логики.
* **Domain Services:** (`availability.go`) — логика, затрагивающая несколько сущностей.


* **Правило:** Никаких тегов `json` или `sql` в сущностях. Только чистый Go.

### 🟡 2. Application Layer (`internal/service/`)

**Сценарии использования (Use Cases).** Отвечает за оркестрацию ("что сделать"), но не за бизнес-правила ("как сделать").

* **Задачи:**
1. Получить входные данные.
2. Загрузить Агрегат из Репозитория.
3. Вызвать метод Агрегата.
4. Сохранить изменения через Репозиторий.
5. Отправить уведомление / событие.


* **Пример:** `create_appointment.go`, `fire_employee.go`.

### 🔴 3. Infrastructure Layer

Реализация технических деталей. Скрыта за интерфейсами.

* **Repository (`internal/repository/`):** Работа с БД (Postgres) и кэшем (Redis). Реализует интерфейсы, определенные в слое Application/Domain.
* *Паттерн:* Transactional Outbox (см. `outbox.go`) для надежной публикации событий.


* **Gateway (`internal/gateway/`):** Клиенты внешних сервисов (SMS, WhatsApp, Payment).
* Используем паттерн **ACL (Anti-Corruption Layer)**, чтобы внешние DTO не проникали в домен.



### 🔵 4. Interface Layer

Входные точки в приложение.

* **Handler (`internal/handler/`):** HTTP контроллеры. Парсят Request, валидируют DTO, вызывают Service, пишут Response.
* **DTO (`internal/dto/`):** Структуры для передачи данных по сети (`request/`, `response/`). Отделяют API контракт от доменных моделей.

---

## 3. Ключевые паттерны и подходы

### Multi-Tenancy (Мульти-аренда)

Система обслуживает множество организаций.

* **OrgContext:** Идентификатор организации извлекается в Middleware (`middleware/org_context.go`) и прокидывается через `context.Context` вглубь приложения.
* **Изоляция данных:** Все репозитории обязаны фильтровать данные по `organization_id`.

### Cross-Context Communication (Модульность)

Модули (например, `booking` и `schedule`) должны быть слабо связаны.

* **Внутри одного процесса:** Используем **Adapters** (`internal/adapter/`).
* *Пример:* `schedule_checker.go` позволяет модулю бронирования проверить расписание, не импортируя напрямую репозитории расписания.


* **Между процессами:** Используем **Events** и **Worker**.

### Authorization Policy

Логика прав доступа отделена от бизнес-логики.

* Используем слой **Policy** (`internal/policy/`), который проверяет, имеет ли `User` право выполнить действие над `Resource`.

### Error Handling

* Ошибки оборачиваются.
* На уровне `Handler` происходит маппинг доменных ошибок (напр. `ErrInsufficientFunds`) в HTTP статусы (402 Payment Required).

---

## 4. Workflows

### Работа с фоновыми задачами

* Используем папку `cmd/worker` для запуска асинхронных процессов.
* Воркены обрабатывают задачи из очереди (Redis/DB) или Outbox-таблицы (для уведомлений и интеграций).

### Миграции

* Все изменения схемы БД хранятся в `migrations/*.sql` и накатываются автоматически или через Make-команды.

---

## 5. Сводка правил разработки (Checklist)

1. [ ] **Logic in Domain:** Если это бизнес-правило ("нельзя уволить владельца"), оно должно быть в `internal/domain`.
2. [ ] **Orchestration in Service:** Если это последовательность шагов ("создать юзера -> отправить email"), это в `internal/service`.
3. [ ] **No Circular Deps:** Следите за импортами. Domain не импортирует Service.
4. [ ] **Use DTOs:** Никогда не возвращайте доменные сущности напрямую в JSON ответах.
5. [ ] **Context Propagation:** Всегда прокидывайте `ctx` первым аргументом.