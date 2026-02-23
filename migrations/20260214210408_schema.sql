-- +goose Up
-- +goose StatementBegin
-- ═══════════════════════════════════════════════════════════════
-- 0. EXTENSIONS (Расширения)
-- ═══════════════════════════════════════════════════════════════
CREATE EXTENSION IF NOT EXISTS btree_gist;


-- ═══════════════════════════════════════════════════════════════
-- 0.1 CUSTOM TYPES (Кастомные типы) - FIX ДЛЯ ОШИБКИ TIMERANGE
-- ═══════════════════════════════════════════════════════════════
-- PostgreSQL не имеет встроенного timerange, создаем его сами:
DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'timerange') THEN
            CREATE TYPE timerange AS RANGE (
                                               subtype = time
                                           );
        END IF;
    END$$;

-- ═══════════════════════════════════════════════════════════════
-- 1. ORGANIZATIONS (Организации)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE organizations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    name VARCHAR(255),
    description TEXT,
    slug VARCHAR(500),
    active BOOLEAN DEFAULT TRUE,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_organization_slug
    UNIQUE (slug)
);

-- ═══════════════════════════════════════════════════════════════
-- 2. PLANS & SUBSCRIPTIONS (Тарифы и Подписки)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(255) NOT NULL UNIQUE,
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

CREATE TABLE plan_capabilities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,

    resource VARCHAR(255) NOT NULL,
    limit_value INT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    UNIQUE (plan_id, resource)
);

CREATE TABLE subscriptions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE RESTRICT,
    plan_id UUID NOT NULL REFERENCES plans(id) ON DELETE RESTRICT,
    status VARCHAR(50) NOT NULL,
    billing_cycle VARCHAR(50) NOT NULL,

    recurring_amount DECIMAL(19,4) NOT NULL,

    current_period_start TIMESTAMPTZ,
    current_period_end TIMESTAMPTZ,
    next_billing_at TIMESTAMPTZ,

    next_retry_at TIMESTAMPTZ,
    retry_count INT DEFAULT 0,

    suspended_at TIMESTAMPTZ,
    canceled_at TIMESTAMPTZ,
    data_delete_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT subscription_per_organization UNIQUE (organization_id),
    CONSTRAINT valid_subscription_status CHECK (
    status IN ('trial', 'active', 'past_due', 'suspended', 'canceled')
    )
);

-- ═══════════════════════════════════════════════════════════════
-- 3. LOCATIONS (Локации и категории)
-- ═══════════════════════════════════════════════════════════════

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

CREATE TABLE locations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES location_categories(id) ON DELETE SET NULL,

    name VARCHAR(255) NOT NULL,
    description TEXT,
    phone VARCHAR(50),
    slug VARCHAR(500),

    city VARCHAR(255),
    address_street VARCHAR(255),
    address_building VARCHAR(50),
    address_details VARCHAR(255),

    longitude DOUBLE PRECISION,
    latitude DOUBLE PRECISION,

    active BOOLEAN DEFAULT true,
    schedule_type VARCHAR(20),
    slot_duration_minutes INTEGER NOT NULL,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT unique_location_slug
    UNIQUE (slug)
);

CREATE TABLE location_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES locations(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT valid_time_range CHECK (start_time < end_time),

    -- Теперь timerange существует и этот код сработает
    CONSTRAINT no_overlapping_slots EXCLUDE USING gist (
    location_id WITH =,
    date WITH =,
    timerange(start_time, end_time) WITH &&
    )
);

-- ═══════════════════════════════════════════════════════════════
-- 4. EMPLOYEES (Сотрудники)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE employees (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    phone VARCHAR(20) NOT NULL,
    password_hash VARCHAR(255),
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    full_name VARCHAR(201) GENERATED ALWAYS AS (first_name || ' ' || COALESCE(last_name, '')) STORED,

    organization_id UUID NOT NULL REFERENCES organizations(id),
    location_id UUID REFERENCES locations(id),

    role VARCHAR(50) NOT NULL,

    can_provide_services BOOLEAN DEFAULT false,
    can_manage_location_schedule BOOLEAN DEFAULT false,

    active BOOLEAN DEFAULT true,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

-- Работник может быть только в одной организации
CREATE UNIQUE INDEX unique_employee_phone
    ON employees(phone) WHERE deleted_at IS NULL;

-- Индекс для владельца
CREATE UNIQUE INDEX idx_one_owner_per_org
    ON employees(organization_id)
    WHERE role = 'owner' AND deleted_at IS NULL;


CREATE TABLE employee_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,

    date DATE NOT NULL,
    type VARCHAR(50) NOT NULL,

    start_time TIME NOT NULL,
    end_time TIME NOT NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    CONSTRAINT valid_time_range_emp CHECK (start_time < end_time),

    CONSTRAINT no_overlapping_slots_emp EXCLUDE USING gist (
    employee_id WITH =,
    date WITH =,
    timerange(start_time, end_time) WITH &&
    )
);

CREATE TABLE employees_work_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id),

    organization_id UUID NOT NULL REFERENCES organizations(id),
    location_id UUID REFERENCES locations(id),

    role VARCHAR(50) NOT NULL,

    started_at TIMESTAMPTZ NOT NULL,
    ended_at TIMESTAMPTZ,

    change_type VARCHAR(50) NOT NULL,
    comment TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT valid_history_period CHECK (ended_at IS NULL OR ended_at > started_at)
);

-- ═══════════════════════════════════════════════════════════════
-- 5. CUSTOMERS (Клиенты)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    phone VARCHAR(20) UNIQUE NOT NULL,

    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100),
    full_name VARCHAR(201) GENERATED ALWAYS AS (first_name || ' ' || COALESCE(last_name, '')) STORED,

    birth_date DATE,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE TABLE customers_base(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    customer_id UUID NOT NULL REFERENCES customers(id) ON DELETE CASCADE,

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

-- ═══════════════════════════════════════════════════════════════
-- 6. SERVICES (Услуги)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE service_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_category_id UUID NOT NULL REFERENCES location_categories(id),
    parent_id UUID REFERENCES service_categories(id),
    name VARCHAR(100) NOT NULL,
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    location_id UUID NOT NULL REFERENCES locations(id),
    category_id UUID NOT NULL REFERENCES service_categories(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    duration_minutes INT NOT NULL,
    color VARCHAR(50) NOT NULL DEFAULT 'gray',
    sort_order INT DEFAULT 0,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    CONSTRAINT services_duration_positive CHECK (duration_minutes > 0),
    CONSTRAINT services_sort_order_non_negative CHECK (sort_order >= 0),
    CONSTRAINT services_name_not_empty CHECK (LENGTH(TRIM(name)) > 0)
);

CREATE UNIQUE INDEX idx_services_unique_name
    ON services(location_id, LOWER(TRIM(name)))
    WHERE deleted_at IS NULL;

CREATE TABLE employee_services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    employee_id UUID NOT NULL REFERENCES employees(id) ON DELETE CASCADE,
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    price DECIMAL(19, 4) NOT NULL,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    UNIQUE (employee_id, service_id)
);

-- ═══════════════════════════════════════════════════════════════
-- 7. APPOINTMENTS (Записи)
-- ═══════════════════════════════════════════════════════════════

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    organization_id UUID NOT NULL REFERENCES organizations(id),

    location_id UUID REFERENCES locations(id) ON DELETE SET NULL,
    service_id UUID REFERENCES services(id) ON DELETE SET NULL,
    employee_id UUID REFERENCES employees(id) ON DELETE SET NULL,
    customer_id UUID REFERENCES customers(id) ON DELETE SET NULL,

    start_at TIMESTAMPTZ NOT NULL,
    end_at TIMESTAMPTZ NOT NULL,

    price DECIMAL(19, 4) NOT NULL,
    duration_minutes INTEGER NOT NULL,

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

    from_status VARCHAR(50),
    to_status VARCHAR(50) NOT NULL,
    payload JSONB,
    changed_by UUID,
    reason TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- ═══════════════════════════════════════════════════════════════
-- 8. NOTIFICATIONS & PARTITIONING (Уведомления)
-- ═══════════════════════════════════════════════════════════════
CREATE TABLE notification_outbox (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY, -- Простой PK

    entity_id TEXT NOT NULL,
    type TEXT NOT NULL, -- '24h_reminder', '1h_reminder', etc.

    payload JSONB NOT NULL,
    scheduled_for TIMESTAMPTZ NOT NULL,
    locked_until TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT now()
);



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- ═══════════════════════════════════════════════════════════════
-- 8. NOTIFICATIONS
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS notification_outbox;

-- ═══════════════════════════════════════════════════════════════
-- 7. APPOINTMENTS
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS appointment_histories;
DROP TABLE IF EXISTS appointments;

-- ═══════════════════════════════════════════════════════════════
-- 6. SERVICES
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS employee_services;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS service_categories;

-- ═══════════════════════════════════════════════════════════════
-- 5. CUSTOMERS
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS customers_notes;
DROP TABLE IF EXISTS customers_base;
DROP TABLE IF EXISTS customers;

-- ═══════════════════════════════════════════════════════════════
-- 4. EMPLOYEES
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS employees_work_history;
DROP TABLE IF EXISTS employee_schedules;
DROP TABLE IF EXISTS employees;

-- ═══════════════════════════════════════════════════════════════
-- 3. LOCATIONS
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS location_schedules;
DROP TABLE IF EXISTS locations;
DROP TABLE IF EXISTS location_categories;

-- ═══════════════════════════════════════════════════════════════
-- 2. PLANS & SUBSCRIPTIONS
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS subscriptions;
DROP TABLE IF EXISTS plan_capabilities;
DROP TABLE IF EXISTS plans;

-- ═══════════════════════════════════════════════════════════════
-- 1. ORGANIZATIONS
-- ═══════════════════════════════════════════════════════════════
DROP TABLE IF EXISTS organizations;

-- ═══════════════════════════════════════════════════════════════
-- 0.1 CUSTOM TYPES
-- ═══════════════════════════════════════════════════════════════
-- Удаляем кастомный тип timerange.
-- Используем CASCADE, чтобы удалить и зависимые от него операторы/функции, если они создались неявно.
DROP TYPE IF EXISTS timerange CASCADE;

-- ═══════════════════════════════════════════════════════════════
-- 0. EXTENSIONS
-- ═══════════════════════════════════════════════════════════════
-- Удаляем расширение.
-- Внимание: если btree_gist используется в других миграциях, эту строку лучше закомментировать.
DROP EXTENSION IF EXISTS btree_gist;

-- +goose StatementEnd
