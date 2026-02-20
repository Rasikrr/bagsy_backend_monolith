# План переноса фичи записи на услугу (Booking) в новую архитектуру

Этот документ описывает стратегию переноса функционала "Bagsy" из ветки `main` в ветку `full-refactor` с учетом новой доменной модели и архитектурных паттернов.

## 1. Сравнение сущностей и моделей

### Доменная модель (Domain)

| `main` (Bagsy) | `full-refactor` (Appointment) | Комментарий |
| :--- | :--- | :--- |
| `ID` (uuid) | `ID` (uuid) | Без изменений |
| `PointCode` (string) | `LocationID` (uuid) | Переход от строковых кодов к UUID локаций |
| `ServiceID` (uuid) | `ServiceID` (uuid) | ID услуги из каталога |
| `ClientPhone` (string) | `CustomerID` (uuid) | Клиенты теперь полноценные сущности `Customer` |
| `MasterPhone` (string) | `EmployeeID` (uuid) | Мастера теперь `Employee` с ролью `Master` |
| `Status` (Status) | `Status` + `StatusHistory` | История статусов теперь внутри агрегата |
| `Price` (decimal) | `Price` (shared.Money) | Использование VO для денег |
| `StartAt`, `EndAt` | `StartAt`, `EndAt` | Без изменений |
| `Comment` | `CustomerComment` | Переименование |
| `RejectReason` | `CancellationReason` | Переименование |
| - | `OrganizationID` (uuid) | Новое обязательное поле (Multi-tenancy) |
| - | `DurationMinutes` | Теперь хранится в записи для истории |

## 2. Архитектурные изменения

В новой архитектуре мы уходим от `internal/services` в сторону `internal/usecases` и DDD.

### Новые компоненты:
1.  **Domain Aggregate**: `Appointment` в `internal/domain/booking/`.
2.  **Use Cases**: `internal/usecases/booking/`.
    - `CreateBookingUseCase`: Создание записи (аналог `bagsy.Create`).
    - `ConfirmBookingUseCase`: Подтверждение записи кодом.
    - `GetAvailableSlotsUseCase`: Получение свободных слотов.
    - `CancelBookingUseCase`: Отмена записи.
3.  **Repositories**: `internal/repositories/booking/`.
    - Реализация `AppointmentRepository` на SQL.

## 3. План реализации по шагам

### Шаг 1: Подготовка Домена (Booking Domain)
*Доменная модель уже частично создана в `internal/domain/booking/appointment.go`.*
- [ ] Доработать `Appointment` агрегат (добавить недостающие методы, если нужно).
- [ ] Описать интерфейс `AppointmentRepository` в доменном слое.

### Шаг 2: Реализация Use Cases (Application Layer)
- [ ] **CreateBookingUseCase**:
    - Поиск или создание `Customer` по телефону (через `Identity` домен).
    - Проверка существования `Service` и `Employee` (через `Catalog` и `Identity`).
    - Расчет `EndAt` на основе длительности из `Catalog`.
    - Генерация кода подтверждения (через `Auth` / `ActionToken`).
    - Сохранение `Appointment` в статусе `Pending`.
    - Отправка уведомления.
- [ ] **ConfirmBookingUseCase**:
    - Проверка кода подтверждения.
    - Перевод `Appointment` в статус `Confirmed` через метод агрегата.
    - Планирование уведомлений-напоминаний.
- [ ] **GetAvailableSlotsUseCase**:
    - Портирование логики `slots.go` из `main`.
    - Использование `location_schedule.go` и `employee_schedule.go` из домена `Schedule` для получения рабочих часов.
    - Пересечение графиков точки и мастера для определения эффективного рабочего времени.
    - Фильтрация по уже существующим `Appointment` в статусах, отличных от `Cancelled`.
    - Учет `shared.Duration` услуги для расчета длины слота.

### Шаг 3: Реализация Инфраструктуры (Infrastructure Layer)
- [ ] Создать миграции для таблицы `appointments` и `appointment_status_history`.
- [ ] Реализовать `SqlAppointmentRepository` в `internal/repositories/booking/`.
- [ ] Реализовать сохранение истории статусов при каждом обновлении агрегата.

### Шаг 4: Адаптация HTTP Ports
- [ ] Создать хендлеры в `internal/ports/http/handlers/booking/`.
- [ ] Описать DTO для запросов и ответов.
- [ ] Зарегистрировать роуты в `server.go`.

## 4. Зависимости и интеграция

Для работы `Booking` понадобятся следующие существующие (или требующие доработки) компоненты:
1.  **Identity**: Получение данных о мастерах и клиентах.
2.  **Catalog**: Получение данных об услугах и ценах.
3.  **Location**: Получение рабочих часов точек.
4.  **Notification**: Отправка SMS/WhatsApp.
5.  **Schedule**: Проверка занятости мастеров и их личных графиков.

## 5. Риски и нюансы
- **Маппинг данных**: Если в БД уже есть данные `Bagsy`, потребуется скрипт миграции для сопоставления `PointCode -> LocationID` и `Phone -> UserID`.
- **Слоты**: Логика генерации слотов самая сложная часть, ее нужно тщательно протестировать новыми unit-тестами.
- **Транзакции**: Создание юзера + создание записи должно быть атомарным.
