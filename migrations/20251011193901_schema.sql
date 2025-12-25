-- +goose Up
-- +goose StatementBegin

-- ============================================
-- СПРАВОЧНИКИ (DICTIONARIES)
-- ============================================

-- CHECKED (RASSUL)
-- Сети заведений
CREATE TABLE IF NOT EXISTS networks (
    code        TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    deleted_at  TIMESTAMPTZ,
    updated_by  TEXT NOT NULL DEFAULT 'system'
);

-- CHECKED (RASSUL)
-- Категории точек (Салон красоты, Барбершоп, СПА и т.д.)
CREATE TABLE IF NOT EXISTS point_categories (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    updated_by  TEXT NOT NULL DEFAULT 'system'
);

-- CHECKED (RASSUL)
-- Категории услуг (Стрижки, Окрашивание, Маникюр и т.д.)
CREATE TABLE IF NOT EXISTS service_categories (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    updated_by  TEXT NOT NULL DEFAULT 'system'
);

-- CHECKED (RASSUL)
-- Подкатегории услуг
CREATE TABLE IF NOT EXISTS service_subcategories (
    id          SERIAL PRIMARY KEY,
    service_category_id INTEGER NOT NULL, -- FK service_categories
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now(),
    updated_by  TEXT NOT NULL DEFAULT 'system'
);

-- ============================================
-- ОСНОВНЫЕ СУЩНОСТИ (CORE ENTITIES)
-- ============================================

-- CHECKED (RASSUL)
-- Точки обслуживания
CREATE TABLE IF NOT EXISTS points (
    code         TEXT PRIMARY KEY,
    name         TEXT NOT NULL,
    description  TEXT,
    network_code TEXT NOT NULL,
    category_id  INTEGER NOT NULL,
    address      JSONB NOT NULL DEFAULT '{}',
    city         TEXT NOT NULL,
    active       BOOLEAN NOT NULL DEFAULT false,
    schedule     JSONB NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ DEFAULT now(),
    deleted_at   TIMESTAMPTZ,
    updated_by   TEXT NOT NULL DEFAULT 'system',
    CONSTRAINT point_code_network_code_unique UNIQUE (code, network_code)
);

-- CHECKED (RASSUL)
-- Пользователи (клиенты, мастера, администраторы)
CREATE TABLE IF NOT EXISTS users (
    phone      TEXT PRIMARY KEY,
    password   TEXT,
    role       TEXT NOT NULL,
    name       TEXT,
    surname    TEXT,
    point_code TEXT,
    network_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT NOT NULL DEFAULT 'system',
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- CONTINUE HERE !! CHECK updated_by field in all tables
-- Услуги (привязаны к конкретной точке)
CREATE TABLE IF NOT EXISTS services (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    point_code       TEXT NOT NULL,
    category_id      INTEGER NOT NULL,
    subcategory_id   INTEGER,
    name             TEXT NOT NULL,
    description      TEXT,
    duration_minutes INTEGER NOT NULL,
    active           BOOLEAN DEFAULT false,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ DEFAULT now(),
    updated_by TEXT NOT NULL DEFAULT 'system',
    CONSTRAINT services_point_code_cat_subcat_unique UNIQUE (point_code, category_id, subcategory_id)
);

-- ============================================
-- СВЯЗУЮЩИЕ ТАБЛИЦЫ (RELATIONS)
-- ============================================

-- Услуги мастеров с ценами
CREATE TABLE IF NOT EXISTS master_services (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_phone TEXT NOT NULL,
    service_id   UUID NOT NULL,
    price        DECIMAL(10,2) NOT NULL,
    active       BOOLEAN DEFAULT false,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ DEFAULT now(),
    updated_by   TEXT NOT NULL DEFAULT 'system',
    UNIQUE(master_phone, service_id)
);

-- Брони (bagsies)
CREATE TABLE IF NOT EXISTS bagsies (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id   UUID NOT NULL,
    point_code   TEXT NOT NULL,
    client_phone TEXT NOT NULL,
    master_phone TEXT NOT NULL,
    status       TEXT NOT NULL,
    price DECIMAL(19, 4) NOT NULL DEFAULT 0.0000,
    start_at     TIMESTAMPTZ NOT NULL,
    end_at       TIMESTAMPTZ NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ DEFAULT now(),
    updated_by   TEXT NOT NULL DEFAULT 'system'
);


-- Форма обратной связи от клиентов
CREATE TABLE IF NOT EXISTS client_forms (
    id          SERIAL PRIMARY KEY,
    first_name  VARCHAR(50) NOT NULL,
    last_name   VARCHAR(50) NOT NULL,
    role        VARCHAR(50) NOT NULL,
    phone       VARCHAR(20) NOT NULL,
    description TEXT NOT NULL,
    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Удаляем в обратном порядке (от зависимых к независимым)
DROP TABLE IF EXISTS client_forms;
DROP TABLE IF EXISTS bagsies;
DROP TABLE IF EXISTS master_services;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS points;
DROP TABLE IF EXISTS service_subcategories;
DROP TABLE IF EXISTS service_categories;
DROP TABLE IF EXISTS point_categories;
DROP TABLE IF EXISTS networks;

-- +goose StatementEnd
