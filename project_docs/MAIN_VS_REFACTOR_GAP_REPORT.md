# Gap Report: main vs full-refactor

## 1. Полный список эндпоинтов на main (37 шт.)

### Auth (`/api/v1/auth`) — 12 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 1 | POST | `/auth/login` | - | Логин по phone + password, возврат JWT пары |
| 2 | POST | `/auth/refresh` | - | Ротация refresh токена |
| 3 | POST | `/auth/logout` | - | Инвалидация refresh токена |
| 4 | GET | `/auth/verify-auth-token/{token}` | - | Проверка одноразового токена (invite/reset) |
| 5 | POST | `/auth/management/register` | - | Шаг 1 — регистрация owner, отправка OTP |
| 6 | POST | `/auth/management/register/confirm` | - | Шаг 2 — подтверждение OTP, создание user+network |
| 7 | POST | `/auth/management/register/resend` | - | Повторная отправка OTP регистрации |
| 8 | POST | `/auth/staff/register` | mgr+ | Шаг 1 — приглашение staff/manager, отправка ссылки |
| 9 | POST | `/auth/staff/register/confirm` | - | Шаг 2 — завершение регистрации по ссылке |
| 10 | POST | `/auth/staff/register/resend` | mgr+ | Повторная отправка ссылки приглашения |
| 11 | POST | `/auth/password/change` | - | Шаг 1 — запрос сброса пароля |
| 12 | POST | `/auth/password/change/confirm` | - | Шаг 2 — подтверждение сброса пароля |

### Users (`/api/v1/users`, `/staff`, `/customers`) — 6 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 13 | GET | `/users/me` | any | Профиль текущего пользователя |
| 14 | PUT | `/users/me` | any | Обновление имени/фамилии |
| 15 | PUT | `/users/me/schedule` | any | Обновление расписания сотрудника |
| 16 | DELETE | `/users/me/avatar` | any | Удаление аватара |
| 17 | GET | `/staff` | mgr+ | Список сотрудников с пагинацией и фильтрами |
| 18 | GET | `/customers` | any | Список клиентов с пагинацией и фильтрами |

### Bagsies/Bookings (`/api/v1/bagsies`) — 6 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 19 | POST | `/bagsies` | - | Создание записи + отправка OTP |
| 20 | POST | `/bagsies/confirm` | - | Подтверждение записи по OTP |
| 21 | POST | `/bagsies/resend` | - | Повтор OTP для записи |
| 22 | POST | `/bagsies/slots` | - | Доступные слоты на 2 недели |
| 23 | POST | `/bagsies/slots/day` | - | Доступные слоты на конкретный день |
| 24 | POST | `/bagsies/master` | staff+ | Создание записи мастером (без OTP) |

### Calendar (`/api/v1/calendar`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 25 | GET | `/calendar` | any | Календарь записей за период (max 35 дней) |

### Points (`/api/v1/points`) — 2 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 26 | POST | `/points` | owner | Создание точки |
| 27 | GET | `/points/{code}` | - | Публичная информация о точке |

### Networks (`/api/v1/networks`) — 2 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 28 | GET | `/networks/{code}` | - | Информация о сети |
| 29 | GET | `/networks/{code}/points` | - | Список точек сети |

### Services (`/api/v1/services`) — 2 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 30 | GET | `/services/{point_code}` | - | Список услуг точки |
| 31 | POST | `/services` | mgr+ | Создание услуги |

### Master Services (`/api/v1/master-services`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 32 | POST | `/master-services` | staff+ | Привязка мастера к услуге + цена |

### Media (`/api/v1/media`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 33 | POST | `/media/upload` | any | Генерация presigned S3 URL для загрузки |

### Point Categories (`/api/v1/point-categories`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 34 | GET | `/point-categories` | - | Список категорий точек |

### Service Categories (`/api/v1/service-categories`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 35 | GET | `/service-categories/{point_code}` | - | Категории услуг точки |

### Forms (`/api/v1/forms`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 36 | POST | `/forms` | - | Заявка на партнёрство |

### Swagger — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 37 | GET | `/swagger/*` | - | Swagger UI (non-prod) |

---

## 2. Эндпоинты на full-refactor (20 шт.)

### Auth (`/api/v1/auth`) — 9 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 1 | POST | `/auth/register` | - | Начало регистрации owner (OTP) |
| 2 | POST | `/auth/register/verify` | - | Подтверждение регистрации (OTP) |
| 3 | POST | `/auth/register/resend` | - | Повтор OTP регистрации |
| 4 | POST | `/auth/login` | - | Логин сотрудника |
| 5 | POST | `/auth/refresh` | - | Ротация токенов |
| 6 | POST | `/auth/logout` | - | Логаут |
| 7 | POST | `/auth/password/reset` | - | Запрос сброса пароля |
| 8 | POST | `/auth/password/reset/confirm` | - | Подтверждение сброса |
| 9 | GET | `/auth/verify/{token}` | - | Проверка action token |

### Employees (`/api/v1/employees`) — 3 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 10 | POST | `/employees/invite` | auth+org | Приглашение сотрудника |
| 11 | POST | `/employees/invite/resend` | auth+org | Повтор приглашения |
| 12 | POST | `/employees/invite/confirm` | - | Завершение регистрации по invite |

### Bookings (`/api/v1/bookings`) — 6 endpoints

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 13 | POST | `/bookings` | - | Создание записи |
| 14 | POST | `/bookings/slots` | - | Доступные слоты |
| 15 | POST | `/bookings/{id}/confirm` | - | Подтверждение OTP |
| 16 | POST | `/bookings/{id}/resend-otp` | - | Повтор OTP |
| 17 | POST | `/bookings/{id}/cancel` | auth+org | Отмена записи |
| 18 | GET | `/bookings/calendar` | auth+org | Календарь записей |

### Locations (`/api/v1/locations`) — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 19 | POST | `/locations` | auth+org | Создание локации |

### Swagger — 1 endpoint

| # | Method | Path | Auth | Описание |
|---|--------|------|------|----------|
| 20 | GET | `/swagger/*` | - | Swagger UI |

---

## 3. Gap Analysis — чего нет в full-refactor

### Критичные (core business logic)

| # | Функционал (main) | Приоритет | Примечание |
|---|-------------------|-----------|------------|
| 1 | **GET `/users/me`** — профиль текущего пользователя | HIGH | Базовый функционал для любого авторизованного пользователя |
| 2 | **PUT `/users/me`** — обновление имени/фамилии | HIGH | CRUD профиля |
| 3 | **PUT `/users/me/schedule`** — обновление расписания | HIGH | Без этого мастер не может настроить своё рабочее время |
| 4 | **GET `/staff`** — список сотрудников | HIGH | Нужен manager/owner для управления персоналом |
| 5 | **GET `/customers`** — список клиентов | HIGH | Нужен для CRM-функционала |
| 6 | **POST `/bagsies/master`** — создание записи мастером | HIGH | Мастер записывает клиента без OTP (основной сценарий для салонов) |
| 7 | **POST `/services`** — создание услуги | HIGH | Без этого нельзя добавлять услуги в каталог |
| 8 | **GET `/services/{point_code}`** — список услуг | HIGH | Публичный просмотр услуг |
| 9 | **POST `/master-services`** — привязка мастера к услуге | HIGH | Без этого мастер не может оказывать услуги |

### Средние (important but not blocking)

| # | Функционал (main) | Приоритет | Примечание |
|---|-------------------|-----------|------------|
| 10 | **POST `/bagsies/slots/day`** — слоты на конкретный день | MEDIUM | В refactor есть slots, но нет day-specific варианта |
| 11 | **GET `/points/{code}`** — публичная страница точки | MEDIUM | Нужен для клиентского приложения |
| 12 | **GET `/networks/{code}`** — информация о сети | MEDIUM | В refactor сети заменены организациями, может быть не нужен |
| 13 | **GET `/networks/{code}/points`** — список точек сети | MEDIUM | Эквивалент: список локаций организации |
| 14 | **POST `/media/upload`** — загрузка файлов | MEDIUM | Нужен для аватаров, фото точек |
| 15 | **DELETE `/users/me/avatar`** — удаление аватара | MEDIUM | Зависит от media |

### Низкие (nice to have / изменившийся контекст)

| # | Функционал (main) | Приоритет | Примечание |
|---|-------------------|-----------|------------|
| 16 | **GET `/point-categories`** — категории точек | LOW | Справочник, легко добавить |
| 17 | **GET `/service-categories/{point_code}`** — категории услуг | LOW | Справочник |
| 18 | **POST `/forms`** — заявка на партнёрство | LOW | Маркетинговый функционал |

### Не нужно переносить (изменения в архитектуре)

| Функционал (main) | Причина |
|-------------------|---------|
| `POST /auth/management/register` + `/confirm` + `/resend` | Заменён на `POST /auth/register` + `/verify` + `/resend` (единый flow для owner) |
| `POST /auth/staff/register` + `/confirm` + `/resend` | Заменён на `POST /employees/invite` + `/confirm` + `/resend` (отдельный контекст) |
| `POST /auth/password/change` + `/confirm` | Заменён на `POST /auth/password/reset` + `/confirm` (переименование) |
| `GET /networks/{code}` | Концепция "сеть" заменена на "организация". Эквивалент — данные из OrgContext |
| `POST /bagsies/confirm` (body: bagsy_id + code) | Заменён на `POST /bookings/{id}/confirm` (id в path) |
| `POST /bagsies/resend` (body: bagsy_id) | Заменён на `POST /bookings/{id}/resend-otp` (id в path) |

---

## 4. Рекомендуемый порядок реализации

```
Волна 1 — каталог и услуги (блокирует остальное):
  [ ] POST /services — создание услуги
  [ ] GET  /services — список услуг локации
  [ ] POST /master-services — привязка мастера к услуге

Волна 2 — профиль и персонал:
  [ ] GET  /users/me — профиль
  [ ] PUT  /users/me — обновление профиля
  [ ] PUT  /users/me/schedule — расписание сотрудника
  [ ] GET  /staff — список сотрудников

Волна 3 — бронирование (расширение):
  [ ] POST /bookings/master — запись мастером без OTP
  [ ] POST /bookings/slots/day — слоты на один день

Волна 4 — клиенты и публичное API:
  [ ] GET  /customers — список клиентов
  [ ] GET  /locations/{id} — публичная страница локации
  [ ] GET  /organizations/{id}/locations — список локаций

Волна 5 — медиа и справочники:
  [ ] POST /media/upload — загрузка файлов
  [ ] DELETE /users/me/avatar
  [ ] GET  /location-categories — справочник
  [ ] GET  /service-categories — справочник
```
