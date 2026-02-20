# Анализ реализации Booking (Bagsy) и гид по рефакторингу

## 1. Что сделано хорошо (Strengths)
1. **Separation of Concerns (Commands):** Использование `CreateBagsyCommand` и `GetAvailableSlotsCommand` разделяет API и бизнес-логику.
2. **Transaction Integrity:** Атомарность создания пользователя и записи.
3. **Pure Slot Logic:** Алгоритм генерации слотов вынесен в функции, что упрощает их тестирование.

## 2. Что требует исправления (Refactoring Targets)
1. **Anemic Domain:** Сущность `Bagsy` не имеет поведения. Все проверки статусов лежат в сервисе.
   - *Refactor:* Перенести логику в методы `Appointment` (`Confirm`, `Cancel`, и т.д.).
2. **Primitive Obsession:** Телефоны и цены представлены строками и числами.
   - *Refactor:* Использовать `shared.Phone`, `shared.Money`.
3. **Implicit Timezones:** Магические преобразования в Almaty Time.
   - *Refactor:* Использовать `TimeZone` из агрегата `Location`.
4. **Fat Orchestration:** Сервис делает слишком много (уведомления, кэш, БД).
   - *Refactor:* Разбить на `UseCases` и внедрить `Domain Events` для уведомлений.

## 3. Маппинг данных (Main -> Refactor)
- `PointCode` (string) -> `LocationID` (uuid)
- `MasterPhone` (string) -> `EmployeeID` (uuid)
- `ClientPhone` (string) -> `CustomerID` (uuid)
- `BagsyStatus` -> `AppointmentStatus` + `StatusHistory`

## 4. План реализации слотов
1. Получить `LocationSchedule`.
2. Получить `EmployeeSchedule`.
3. Найти пересечение (Effective Working Hours).
4. Вычесть существующие `Appointments` и `Breaks`.
5. Разбить остаток на интервалы согласно `Service.Duration`.
