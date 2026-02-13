--------------------------------------
        -- Employees --
--------------------------------------

CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255),  -- NULL до принятия инвайта
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),

    organization_id UUID NOT NULL REFERENCES organizations(id),
    location_id UUID NOT NULL REFERENCES locations(id),

    role VARCHAR(50) NOT NULL,  -- owner, manager, staff
    can_provide_services BOOLEAN DEFAULT false,
    can_manage_location_schedule BOOLEAN DEFAULT false,
    -- TODO: add other attributes --

    active BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    -- Сотрудник может быть только в одной организации
    CONSTRAINT unique_active_employee_phone
       UNIQUE (phone) WHERE active = true
);

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE employee_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL, -- work, rest --
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT valid_time_range CHECK (start_time < end_time),

    -- Запрет пересечений
    CONSTRAINT no_overlapping_slots EXCLUDE USING gist (
        employee_id WITH =,
        date WITH =,
    timerange(start_time, end_time) WITH &&
    )
);


CREATE TABLE employees_work_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),
    organization_id UUID FOREIGN KEY REFERENCES organizations(id),
    role VARCHAR(50) NOT NULL, -- Важно хранить, кем он был
    joined_at TIMESTAMPTZ NOT NULL,  -- Когда принял инвайт или сменил должность в этой же организации
    fired_at TIMESTAMPTZ, -- Когда увелен с последнего метса работы
    fire_reason TEXT
)



--------------------------------------
        -- Customers --
--------------------------------------


CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(20) UNIQUE NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    birth_date DATE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,
);

 -- База клиентов, куда салоны могут заполнять информацию о клиенте --
 -- Информация о клиенте распространяется на всю сеть (не на отдельные точки) --
CREATE TABLE customers_base(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID FOREIGN KEY REFERENCES customers(id),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    birth_date DATE,
    gender VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    CONSTRAINT unique_customer_base_organization
        UNIQUE (organization_id, customer_id)
);

CREATE TABLE customers_notes(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_base_id UUID NOT NULL REFERENCES customers_base(id),
    author_id UUID NOT NULL REFERENCES employees(id),
    note TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);




--------------------------------------
        -- Organizations --
--------------------------------------

CREATE TABLE organizations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID FOREIGN KEY REFERENCES employees(id),
    name VARCHAR(255),
    description TEXT,
    slug VARCHAR(500),
    active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_organization_slug
        UNIQUE (slug);
);


-- ═══════════════════════════════════════════════════════════════
-- ПЛАНЫ (тарифы)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,            -- 'Solo', 'Business', 'Enterprise'
    description TEXT,

    -- Цены
    price_monthly DECIMAL(19,4) NOT NULL,
    price_annual DECIMAL(19,4) NOT NULL,

    -- Мета
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

-- ═══════════════════════════════════════════════════════════════
-- ФИЧИ ПЛАНОВ
-- ═══════════════════════════════════════════════════════════════
-- Тут также будут лимиты:
-- Лимиты
--     max_locations 1,      -- -1 = безлимит
--     max_employees 1,      -- -1 = безлимит
--     max_services  1       -- -1 = безлимит
--     analytics_enabled NULL -- если строки нет - недоступно

CREATE TABLE plan_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,

    resource VARCHAR(255) NOT NULL,
    -- Лимит для фичи (опционально)
    limit_value INT,  -- NULL = безлимит, число = лимит

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    UNIQUE (plan_id, resource)
);
-- ═══════════════════════════════════════════════════════════════
-- ПОДПИСКИ (ссылается на plan)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    plan_id UUID NOT NULL REFERENCES plans(id),
    billing_cycle VARCHAR(100) NOT NULL,

    -- !!! SNAPSHOT ЦЕНЫ !!!
    -- Важно сохранять, по какой цене клиент подписался, даже без валюты.
    -- Если завтра тариф станет дороже, этот клиент продолжит платить эту сумму.
    recurring_amount DECIMAL(19,2) NOT NULL,

    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    trial_ends_at TIMESTAMPTZ,
    next_billing_at TIMESTAMPTZ,
    suspended_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT subscription_per_organization UNIQUE (organization_id)
);

-- ═══════════════════════════════════════════════════════════════
-- ПЛАТЕЖИ
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE subscription_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    subscription_id UUID NOT NULL REFERENCES subscriptions(id),
    amount DECIMAL(19,4) NOT NULL,
    status VARCHAR(100) NOT NULL,-- pending, success, failed, refunded --
    payment_provider VARCHAR(255),
    external_payment_id VARCHAR(255),
    paid_at TIMESTAMPTZ,
    fail_reason TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);
-------



--------------------------------------
            -- Locations --
--------------------------------------

CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),
    category_id UUID NOT NULL REFERENCES location_categories(id),

    name VARCHAR(255) NOT NULL,
    description TEXT,
    slug VARCHAR(500),
    city VARCHAR(255),
    address_street VARCHAR(255),
    address_building VARCHAR(50),
    address_details VARCHAR(255),

    longitude DOUBLE PRECISION,
    latitude DOUBLE PRECISION,

    active BOOLEAN DEFAULT true,
    schedule_type VARCHAR(20), --fixed, mixed--
    slot_duration_minutes INTEGER NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_location_slug
        UNIQUE (slug)
);


CREATE TABLE location_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug VARCHAR(500) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT unique_location_categories_slug
        UNIQUE (slug)
);

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE TABLE location_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL, -- work, rest --
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT valid_time_range CHECK (start_time < end_time),

    -- Запрет пересечений
    CONSTRAINT no_overlapping_slots EXCLUDE USING gist (
        location_id WITH =,
        date WITH =,
        timerange(start_time, end_time) WITH &&
    )
);



--------------------------------------
        -- Services --
--------------------------------------

CREATE TABLE service_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_category_id UUID NOT NULL REFERENCES location_categories(id),
    parent_id UUID REFERENCES service_categories(id),  -- Для подкатегорий
    name VARCHAR(100) NOT NULL,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    -- TODO: maybe add slug --
);


CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES locations(id),
    category_id UUID NOT NULL REFERENCES service_categories(id),  -- NULL = без категории
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INT NOT NULL,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT true, -- Временно скрыть
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ -- Удалить навсегда (soft для отчетов, аналитики)

    CONSTRAINT services_duration_positive CHECK (duration_minutes > 0),
    CONSTRAINT services_sort_order_non_negative CHECK (sort_order >= 0),
    CONSTRAINT services_name_not_empty CHECK (LENGTH(TRIM(name)) > 0)
);

-- Уникальность имени услуги
CREATE UNIQUE INDEX idx_services_unique_name
    ON services(location_id, LOWER(TRIM(name)))
    WHERE deleted_at IS NULL;


-- Связь мастер-услуга (с индивидуальной ценой)
CREATE TABLE employee_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    price DECIMAL(19, 4) NOT NULL,                -- Цена этого мастера
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    UNIQUE (employee_id, service_id)
);

--------------------------------------
        -- Appointment --
--------------------------------------

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),

    -- ССЫЛКИ (Nullable для Hard Delete защиты)
    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    service_id UUID REFERENCES services(id) ON DELETE SET NULL,
    employee_id UUID REFERENCES employees(id) ON DELETE SET NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,

    -- Начало и конец записи
    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,

    -- === SNAPSHOTS (ИСТОРИЯ) ===

    -- 1. Финансы (Критично)
    price DECIMAL(19, 4) NOT NULL,
    -- Если добавите себестоимость или скидку, их snapshot тоже сюда:
    -- discount_amount DECIMAL(19, 4) DEFAULT 0,

    duration_minutes INTEGER NOT NULL, -- Snapshot длительности

    -- СТАТУСЫ И МЕТА
    status VARCHAR(50) NOT NULL,
    customer_comment TEXT,

    cancelled_by UUID,
    cancellation_reason VARCHAR(500),

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ
);


CREATE TABLE appointment_histories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    appointment_id UUID NOT NULL REFERENCES appointments(id) ON DELETE CASCADE,

    from_status VARCHAR(50), -- pending, confirmed, in_progress, completed, cancelled
    to_status VARCHAR(50) NOT NULL, -- pending, confirmed, in_progress, completed, cancelled
    payload JSONB,
    changed_by UUID,  -- employee_id
    reason TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW()
);


--------------------------------------
    -- Notification outbox --
--------------------------------------

CREATE TABLE notification_outbox (
    id BIGINT GENERATED ALWAYS AS IDENTITY,

    -- ОБЯЗАТЕЛЬНО ДОБАВИТЬ: ID сущности (Order ID, Appointment ID)
    entity_id TEXT NOT NULL,
    type TEXT NOT NULL, -- '24h_reminder', '1h_reminder', etc.

    payload JSONB NOT NULL,
    scheduled_for TIMESTAMPTZ NOT NULL,
    locked_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (id, scheduled_for)
) PARTITION BY RANGE (scheduled_for);

-- Пример запроса для воркера уведомлений
-- WITH batch AS (
--     SELECT id, scheduled_for
--     FROM notification_outbox
--     WHERE scheduled_for <= NOW()
--       AND (locked_until IS NULL OR locked_until < NOW()) -- Берем новые или "протухшие"
--     ORDER BY scheduled_for ASC
--     LIMIT 50 -- Размер пачки
--     FOR UPDATE SKIP LOCKED -- Пропускаем то, что прямо сейчас лочат другие
--             )
-- UPDATE notification_outbox q
-- SET locked_until = NOW() + INTERVAL '2 minutes' -- Время на попытку отправки
-- FROM batch b
-- WHERE q.id = b.id
--   AND q.scheduled_for = b.scheduled_for -- ВАЖНО для partition pruning!
--     RETURNING q.id, q.scheduled_for, q.type, q.payload;
-- Удаление (Cleanup)
-- Здесь есть нюанс. Нам нужно удалить записи по составному ключу (id, scheduled_for).
-- Простой DELETE WHERE id IN (...) будет медленным, так как Postgres будет сканировать все партиции.

-- func (w *Worker) deleteTasks(ctx context.Context, ids []int64, dates []time.Time) error {
-- 	// Мы используем UNNEST, чтобы удалить пачку за один запрос
-- 	// и при этом использовать индекс по (id, scheduled_for)
--
-- 	query := `
-- 		DELETE FROM notification_outbox
-- 		WHERE (id, scheduled_for) IN (
-- 			SELECT * FROM UNNEST($1::bigint[], $2::timestamptz[])
-- 		)
-- 	`
--
-- 	_, err := w.db.ExecContext(ctx, query, pq.Array(ids), pq.Array(dates))
-- 	return err
-- }


SELECT partman.create_parent(
    -- 1. Имя нашей таблицы
    p_parent_table => 'public.notification_outbox',

    -- 2. Поле, по которому режем (обязательно должно быть в PK)
    p_control => 'scheduled_for',

    -- 3. Тип партиционирования (native - это современный декларативный стиль Postgres)
    p_type => 'native',

    -- 4. Интервал: 1 день (оптимально для очистки)
    p_interval => '1 day',

    -- 5. Самое важное для тебя: PREMAKE
    -- Сколько таблиц держать готовыми ВПЕРЕД.
    -- Ставим 30 (месяц). Это покрывает 90% записей.
    -- Всё, что дальше 30 дней, упадет в _default таблицу.
    p_premake => 30
);


UPDATE partman.part_config
SET
    -- Удалять партиции, где все записи старше 7 дней
    retention = '7 days',

    -- false = DROP TABLE (физическое удаление, освобождает место сразу)
    -- true = DETACH TABLE (просто отцепляет, таблица остается в базе)
    retention_keep_table = false,

    -- Включаем автоматическое перемещение данных из DEFAULT партиции
    -- в правильную партицию, когда она будет создана.
    infinite_time_partitions = true
WHERE parent_table = 'public.notification_outbox';