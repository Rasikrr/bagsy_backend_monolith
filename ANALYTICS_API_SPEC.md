# Analytics API — ТЗ для бэкенда

Спецификация эндпоинтов для раздела «Аналитика» (`/analytics` в ЛК).

> На фронте всё уже реализовано на моках с теми же контрактами. Когда бэк готов — фронт подменяет тела методов в `src/shared/services/analytics-service.ts` на `apiClient.get(...)` и удаляет `src/shared/mocks/analytics.ts`. **Никакие фронтовые компоненты, хуки, типы менять не нужно** — формат ответов должен полностью соответствовать описанному ниже.

---

## 0. Общие правила

### Base URL и авторизация

- Base: `NEXT_PUBLIC_API_URL` (например `https://api.bagsy.kz`)
- Prefix: `/api/v1/`
- Auth: `Authorization: Bearer <access_token>` для всех аналитических эндпоинтов

### Параметры периода

Все эндпоинты аналитики принимают одинаковый набор query-параметров:

| Параметр      | Тип    | Обяз. | Описание                                                        |
| ------------- | ------ | ----- | --------------------------------------------------------------- |
| `from`        | string | ✓     | Дата начала текущего периода (YYYY-MM-DD)                       |
| `to`          | string | ✓     | Дата конца текущего периода включительно (YYYY-MM-DD)           |
| `location_id` | UUID   | —     | Фильтр по локации. Если не передан — все доступные пользователю |

Без `location_id`:

- Для **Manager** возвращаются данные только по его локации (бэк определяет по токену)
- Для **Network Owner** — агрегат по всем локациям организации
- Для **Point Owner** — единственная локация

### Правила расчёта периода сравнения (delta vs prev)

Бэк **сам** вычисляет предыдущий период по `from/to`:

| Длительность текущего периода         | Сравнение с                                  |
| ------------------------------------- | -------------------------------------------- |
| 1 день (`from == to`)                 | Вчера (`from - 1 day`)                       |
| Месяц с 1-го числа по сегодня (MTD)\* | MTD прошлого месяца (1-е прошлого — та же дата) |
| Любая другая длительность N           | Equal-length: N дней до `from`               |

\* Определяется так: `from` = 1-е число месяца И (`to` — последний день того же месяца ИЛИ `to` — сегодня).

Пример equal-length: `from=2026-05-10, to=2026-05-16` (7 дней) → сравнение `2026-05-03..2026-05-09`.

Пример MTD: `from=2026-05-01, to=2026-05-24` → сравнение `2026-04-01..2026-04-24`.

Если в прошлом месяце нет такой даты (например 31-го), берём последний день прошлого месяца.

### Общие типы ответов

```ts
// Значение KPI: текущее, предыдущее, дельта в процентах
type KpiValue = {
  value: number;
  prev: number;
  delta_percent: number | null; // null если prev == 0 (деление на ноль)
};

// Точка временного ряда (день)
type DailyPoint = {
  date: string;     // YYYY-MM-DD
  value: number;    // текущий период
  prev_value: number; // тот же индекс в прошлом периоде
};

// Топ-N сущности (мастер / услуга)
type TopItem = {
  id: string;       // UUID
  name: string;
  revenue: number;  // выручка за период
  share: number;    // доля 0..1 от общей выручки
};

// Этап воронки записей
type FunnelStage = {
  key: "created" | "confirmed" | "completed";
  count: number;
  conversion: number; // 0..1, конверсия от предыдущего этапа; для created = 1
};

// Ячейка heatmap (день недели × час)
type HeatmapCell = {
  weekday: number;  // 0..6, Пн = 0
  hour: number;     // 0..23
  value: number;    // 0..1 — нагрузка
};

// Авто-инсайт (стабильный ключ + параметры для i18n)
type Insight = {
  key: string;      // напр. "saturdayLoad", "revenueDrop"
  level: "info" | "warning" | "success";
  params?: Record<string, string | number>;
};
```

Стабильные ключи инсайтов (фронт мапит на переводы):

| Key                | Уровень   | Параметры          | Условие                                                |
| ------------------ | --------- | ------------------ | ------------------------------------------------------ |
| `saturdayLoad`     | `info`    | `{ percent }`      | Загрузка субботы > 90% от max                          |
| `revenueDrop`      | `warning` | `{ percent }`      | Падение выручки vs пред. период > 5%                   |
| `topServiceShare`  | `success` | `{ name, percent }`| Доля топ-услуги в выручке > 25%                        |
| `retentionFirst`   | `info`    | —                  | Retention после 1-го визита < 70% (на странице clients)|

> Бэк может вернуть пустой массив инсайтов. Фронт умеет это обрабатывать.

### Коды ошибок

```
401 — отсутствует или истёк токен (фронт сам обновит)
403 — нет прав (например Staff пытается GET /staff)
404 — ресурс не найден (например employee_id не существует)
422 — невалидные параметры (from > to, кривой формат даты)
500 — внутренняя ошибка
```

### Кэширование

Желательно `Cache-Control: private, max-age=60` для всех аналитических эндпоинтов — на фронте `staleTime=5min`, бэку достаточно секундной свежести.

---

## 1. `GET /api/v1/analytics/overview` 🔒

Сводка для главной страницы аналитики (Manager/Owner).

### Доступ

- **Owner**, **Manager** — да
- **Staff** → `403`

### Query

```
from, to, location_id?
```

### Response 200

```ts
{
  kpi: {
    revenue: KpiValue;              // выручка в тенге
    bookings: KpiValue;             // количество завершённых записей
    clients: KpiValue;              // уникальных клиентов (по phone)
    avg_check: KpiValue;            // средний чек = revenue / bookings
    load_percent: KpiValue;         // загрузка %: фактические часы записей / часы расписания
    cancellation_percent: KpiValue; // % отменённых от созданных
  };
  revenue_by_day: DailyPoint[];     // выручка по дням периода
  top_employees: TopItem[];         // топ-5 мастеров по выручке
  top_services: TopItem[];          // топ-5 услуг по выручке
  funnel: FunnelStage[];            // [created, confirmed, completed]
  heatmap: HeatmapCell[];           // 7 weekdays × N часов (только рабочие)
  insights: Insight[];              // авто-инсайты
}
```

### Особенности расчёта

- `clients.value` = `COUNT(DISTINCT phone) WHERE appointment IN [from..to] AND status = 'completed'`
- `load_percent.value` = `SUM(duration) / SUM(schedule_hours) * 100` за период
- Heatmap: только рабочие часы локации (например 9..21). Значение нормировано: `cell.value = bookings_in_cell / max(bookings_in_any_cell)`
- Top-5: если мастеров/услуг меньше 5 — возвращаем сколько есть
- `funnel.completed.conversion` = `completed / confirmed`; для `created` всегда 1

### Пример response

```json
{
  "kpi": {
    "revenue": { "value": 1245000, "prev": 1110000, "delta_percent": 12.16 },
    "bookings": { "value": 156, "prev": 144, "delta_percent": 8.33 },
    "clients": { "value": 89, "prev": 77, "delta_percent": 15.58 },
    "avg_check": { "value": 7981, "prev": 7708, "delta_percent": 3.54 },
    "load_percent": { "value": 72, "prev": 68, "delta_percent": 5.88 },
    "cancellation_percent": { "value": 8, "prev": 11, "delta_percent": -27.27 }
  },
  "revenue_by_day": [
    { "date": "2026-05-01", "value": 42000, "prev_value": 38000 },
    { "date": "2026-05-02", "value": 38500, "prev_value": 41000 }
  ],
  "top_employees": [
    { "id": "uuid-1", "name": "Алия", "revenue": 385000, "share": 0.31 }
  ],
  "top_services": [
    { "id": "uuid-s1", "name": "Окрашивание", "revenue": 420000, "share": 0.34 }
  ],
  "funnel": [
    { "key": "created", "count": 178, "conversion": 1 },
    { "key": "confirmed", "count": 165, "conversion": 0.927 },
    { "key": "completed", "count": 156, "conversion": 0.945 }
  ],
  "heatmap": [
    { "weekday": 0, "hour": 9, "value": 0.35 },
    { "weekday": 5, "hour": 17, "value": 0.95 }
  ],
  "insights": [
    {
      "key": "saturdayLoad",
      "level": "info",
      "params": { "percent": 95 }
    }
  ]
}
```

---

## 2. `GET /api/v1/analytics/me` 🔒

Личная аналитика текущего сотрудника. Доступно всем ролям (мастер видит только свою).

### Доступ

- Любой авторизованный — да (бэк фильтрует по `employee_id` из токена)

### Query

```
from, to
```

(`location_id` не нужен — у мастера всегда одна)

### Response 200

```ts
{
  kpi: {
    revenue: KpiValue;
    bookings: KpiValue;
    clients: KpiValue;        // уникальные клиенты этого мастера
    avg_check: KpiValue;
    load_percent: KpiValue;
    cancellation_percent: KpiValue;
  };
  revenue_by_day: DailyPoint[];
  top_services: TopItem[];    // топ-5 услуг этого мастера
  heatmap: HeatmapCell[];     // только часы записей этого мастера
  clients_breakdown: {
    new: number;              // клиенты впервые пришли в этом периоде
    returning: number;        // повторно пришедшие
  };
}
```

### Особенности

- `new` клиент: его первый визит когда-либо у этого мастера попадает в `[from..to]`
- `returning` клиент: первый визит до `from`, но в `[from..to]` тоже был

---

## 3. `GET /api/v1/analytics/staff` 🔒

Отчёт по всем мастерам локации (таблица).

### Доступ

- **Owner**, **Manager** — да
- **Staff** → `403`

### Query

```
from, to, location_id?
```

### Response 200

```ts
{
  rows: Array<{
    employee_id: string;        // UUID
    full_name: string;
    revenue: number;
    bookings: number;
    avg_check: number;
    load_percent: number;       // 0..100
    cancellations: {
      count: number;
      percent: number;          // 0..100
    };
    rating: number | null;      // 1..5, null если не реализовано
  }>;
  weekday_load: Array<{
    employee_id: string;
    weekday: number;            // 0..6
    value: number;              // 0..1
  }>;
  insights: Insight[];
}
```

### Особенности

- `rows` отсортирован по `revenue DESC`
- `weekday_load` — flat-массив, фронт сам группирует в матрицу мастер×день
- Если мастер не работал ни одного дня в периоде — он всё равно в `rows` со всеми нулями (нужно показать что он есть)

---

## 4. `GET /api/v1/analytics/staff/{employee_id}` 🔒

Drill-down по конкретному мастеру.

### Доступ

- **Owner**, **Manager** — да
- **Staff** — только если `employee_id` совпадает с собственным из токена (иначе `403`)

### Path params

```
employee_id: UUID
```

### Query

```
from, to
```

### Response 200

```ts
{
  employee: {
    id: string;
    full_name: string;
    avatar_url?: string;
  };
  kpi: {
    revenue: KpiValue;
    bookings: KpiValue;
    clients: KpiValue;
    avg_check: KpiValue;
    load_percent: KpiValue;
    cancellation_percent: KpiValue;
  };
  revenue_by_day: DailyPoint[];
  top_services: TopItem[];
  hourly_load: Array<{
    hour: number;       // 0..23 (но обычно 9..21)
    value: number;      // 0..1
  }>;
  clients_breakdown: {
    new: number;
    returning: number;
  };
}
```

### Ошибки

- `404` если `employee_id` не найден в организации

---

## 5. `GET /api/v1/analytics/locations/{location_id}` 🔒

Drill-down по конкретной локации. Только для **Network Owner**.

### Доступ

- **Owner** && `subscription.plan == "network"` — да
- Остальные → `403`

### Path params

```
location_id: UUID
```

### Query

```
from, to
```

### Response 200

Идентичен `GET /analytics/overview` (та же структура), но скоупом одна локация.

### Ошибки

- `404` если `location_id` не принадлежит организации пользователя
- `403` если у пользователя не Network-план

---

## 6. `GET /api/v1/analytics/finance` 🔒

Финансовый отчёт.

### Доступ

- **Owner**, **Manager** — да
- **Staff** → `403`

### Query

```
from, to, location_id?
```

### Response 200

```ts
{
  revenue: {
    services: number;  // выручка от услуг
    products: number;  // выручка от товаров (пока 0, до релиза модуля)
    total: number;     // services + products
  };
  payroll: Array<{
    employee_id: string;
    full_name: string;
    commission_percent: number;  // 0..100
    amount: number;              // сумма к выплате
  }>;
  payroll_total: number;
  gross_profit: number;          // total - payroll_total
  margin_percent: number;        // gross_profit / total * 100
}
```

### Особенности

- `payroll[i].amount` = `revenue_от_услуг_мастера * commission_percent / 100`
- Если `commission_percent` не задан для мастера — берём дефолт `0` (UI покажет 0)

### Требуется на стороне БД

Добавить поле `commission_percent` (int 0..100) в `employees`:

```sql
ALTER TABLE employees
ADD COLUMN commission_percent INT NOT NULL DEFAULT 0
CHECK (commission_percent BETWEEN 0 AND 100);
```

UI редактирования этого поля делается отдельной задачей — на текущей итерации фронт просто отображает значение.

---

## 7. `GET /api/v1/analytics/clients` 🔒

Аналитика клиентов (помечено Beta).

### Доступ

- **Owner**, **Manager** — да
- **Staff** → `403`

### Query

```
from, to, location_id?
```

### Response 200

```ts
{
  kpi: {
    total: KpiValue;       // всего уникальных клиентов в базе
    new: KpiValue;         // новые в периоде
    returning: KpiValue;   // вернувшиеся в периоде
    lost: KpiValue;        // не приходят > N дней
  };
  segments: Array<{
    key: "new" | "growing" | "regular" | "vip" | "sleeping" | "lost";
    count: number;
    share: number;         // 0..1
  }>;
  retention: {
    after_1: number;       // 0..1 — % вернувшихся после 1-го визита
    after_2: number;       // после 2-го
    after_3: number;       // после 3-го
  };
  cohorts: Array<{
    month: string;         // "Окт", "Ноя" — или ISO "2026-05" (на ваш выбор; фронт принимает строку)
    new_count: number;
    active_percent: number; // 0..1 — сколько из новых того месяца ещё активны
  }>;
}
```

### Определение сегментов (можно тюнить)

| Сегмент    | Условие                                                     |
| ---------- | ----------------------------------------------------------- |
| `new`      | Первый визит когда-либо в последние 30 дней                  |
| `growing`  | 2-3 визита всего, последний < 60 дней назад                 |
| `regular`  | 4+ визитов, последний < 60 дней назад                       |
| `vip`      | 4+ визитов И средний чек выше 75-й перцентили организации   |
| `sleeping` | Последний визит 60..180 дней назад                          |
| `lost`     | Последний визит > 180 дней назад                            |

Точные пороги обсудим — можем хранить их как настройки организации.

---

## 8. Чеклист реализации

- [ ] Создать таблицу/материализованные view для агрегатов (revenue_by_day, hourly_load — можно считать на лету при <1М записей за период)
- [ ] Добавить поле `employees.commission_percent`
- [ ] Реализовать middleware проверки роли+плана (для `/staff`, `/finance`, `/clients` — нужен `manager_scope`, для `/locations/{id}` — `network_owner`)
- [ ] 7 эндпоинтов: `/overview`, `/me`, `/staff`, `/staff/{id}`, `/locations/{id}`, `/finance`, `/clients`
- [ ] Юнит-тесты на расчёт периода сравнения (today/MTD/equal-length с граничными случаями: 31-е → февраль и т.п.)
- [ ] Юнит-тесты на расчёт `delta_percent` (особенно prev=0 → null)
- [ ] Интеграционные тесты на каждый эндпоинт (happy path + 403 + 404 + 422)
- [ ] OpenAPI описание добавить в swagger
- [ ] Performance budget: каждый эндпоинт < 300мс для периода 90 дней при 5000 записей

---

## 9. Контакты

- Frontend: типы — `src/shared/types/analytics.ts`, моки — `src/shared/mocks/analytics.ts`
- При вопросах: правки в этот файл + согласовать с фронт-командой перед изменением контракта
