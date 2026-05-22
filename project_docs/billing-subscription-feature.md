# Billing & Subscription Feature

## Что реализовано

### Domain Layer

- **`plan_code.go`** — добавлен метод `TrialDays()`: Solo = 60 дней (2 месяца), Point/Network = 30 дней (1 месяц). Соответствует дипломному документу (Chapter 5, Table 5.1).
- **`register_owner.go`** — при регистрации trial создаётся с длительностью зависящей от выбранного плана (`reg.PlanCode.TrialDays()` вместо хардкоженных 60 дней).

### Use Cases (`internal/usecases/billing/`)

| Файл | Метод | Описание |
|------|-------|----------|
| `usecase.go` | — | DI-контейнер: `subscriptionRepo`, `planRepo`, `policy` |
| `get_subscription.go` | `GetSubscription` | Возвращает подписку организации + связанный план |
| `activate_subscription.go` | `Activate` | Переход trial/past_due/suspended → active. Принимает `cycle` (monthly/annual), цена берётся из плана |
| `cancel_subscription.go` | `RequestCancellation` | active → cancel_at_period_end. Подписка остаётся active до конца оплаченного периода |
| `undo_cancellation.go` | `UndoCancellation` | Отмена запроса на отмену (пока период не истёк) |
| `list_plans.go` | `ListPlans` | Список всех активных планов (публичный) |

### Policy (`internal/usecases/policy/`)

- **`CanManageSubscription`** — только Owner может активировать/отменять подписку. Не проверяет `CanOperate()`, чтобы позволить активацию даже из `suspended`.

### Plan Repository

- **`FindAllActive()`** — загружает все активные планы с capabilities, отсортированные по `sort_order`.

### HTTP Handlers (`internal/ports/http/handlers/billing/`)

| Метод | Endpoint | Auth | Описание |
|-------|----------|------|----------|
| GET | `/api/v1/plans` | Публичный | Список тарифов с ценами и лимитами |
| GET | `/api/v1/subscription` | JWT + OrgContext | Текущая подписка организации с планом |
| POST | `/api/v1/subscription/activate` | JWT + OrgContext (Owner) | Активация с выбором цикла |
| POST | `/api/v1/subscription/cancel` | JWT + OrgContext (Owner) | Запрос отмены в конце периода |
| POST | `/api/v1/subscription/undo-cancel` | JWT + OrgContext (Owner) | Отмена запроса на отмену |

Все эндпоинты имеют swagger-аннотации (tag: `billing`).

### Error Mapping

```
subscription_not_found       → 404
subscription_suspended       → 403
subscription_already_active  → 409
invalid_status_transition    → 422
invalid_billing_cycle        → 400
not_active_for_cancellation  → 422
cancellation_already_requested → 409
no_cancellation_to_undo      → 422
plan_not_found               → 404
permission_denied            → 403
```

### Worker (`SubscriptionStatusJob`)

Зарегистрирован в `app.go`. Расписание через env `subscription_worker_schedule`.

Обрабатывает автоматические переходы:
- trial (период истёк) → `past_due`
- active (период истёк, `cancel_at_period_end=false`) → `past_due`
- active (период истёк, `cancel_at_period_end=true`) → `canceled`
- past_due (retry count < 3) → schedule next retry
- past_due (retry count >= 3) → `suspended`
- suspended (90 дней) → `canceled`

### Wiring

- **`app.go`**: `billingUseCase` создаётся и передаётся в `http.NewServer`. Worker зарегистрирован в `initJobs`.
- **`server.go`**: billing handler подключен к роутеру.
- **`envs.go`**: добавлена константа `SubscriptionWorkerSchedule`.

## Жизненный цикл подписки

```
[Регистрация]
     │
     ▼
   trial (Solo=60d, Point/Network=30d)
     │
     ├─ Оплатил досрочно ──► active (monthly/annual)
     │                          │
     │                          ├─ Автопродление ─► active (новый период)
     │                          ├─ RequestCancellation ─► active (cancel_at_period_end=true)
     │                          │                            │
     │                          │                            ├─ UndoCancellation ─► active (cancel_at_period_end=false)
     │                          │                            └─ Период истёк ─► canceled
     │                          │
     │                          └─ Период истёк, не оплатил ─► past_due
     │
     └─ Период истёк ──► past_due
                            │
                            ├─ Оплатил ──► active
                            ├─ Retry (до 3 попыток, каждые 3 дня)
                            └─ Все retry исчерпаны ──► suspended (read-only)
                                                          │
                                                          ├─ Оплатил ──► active
                                                          └─ 90 дней ──► canceled
                                                                           │
                                                                           └─ 90 дней ──► data deletion (TODO)
```

## Что НЕ реализовано

- Payment gateway (Stripe/Kaspi) — activate пока вызывается вручную без реальной оплаты
- Webhook для подтверждения оплаты
- Upgrade/downgrade между тарифами
- Admin-эндпоинты для ручного управления подписками
- Data deletion при `canceled` (в worker стоит TODO)
