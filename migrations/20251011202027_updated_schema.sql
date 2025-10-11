-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS point_categories (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE points
    ADD COLUMN IF NOT EXISTS category_id TEXT,
    ADD COLUMN IF NOT EXISTS active      BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS deleted_at  TIMESTAMPTZ;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS service_categories (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS service_subcategories (
    id          BIGSERIAL PRIMARY KEY,
    category_id TEXT,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS services (
    id               BIGSERIAL PRIMARY KEY,
    point_code       TEXT,
    name             TEXT NOT NULL,
    description      TEXT,
    category_id      TEXT,
    subcategory_id   TEXT,
    duration_minutes TEXT,
    active           BOOLEAN DEFAULT TRUE,
    created_at       TIMESTAMPTZ DEFAULT now(),
    updated_at       TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS master_services (
    id               BIGSERIAL PRIMARY KEY,
    master_phone     TEXT,
    service_id       TEXT,
    price            TEXT,
    duration_minutes TEXT,
    active           BOOLEAN DEFAULT TRUE,
    created_at       TIMESTAMPTZ DEFAULT now(),
    updated_at       TIMESTAMPTZ DEFAULT now()
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS point_code TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE bagsies
    ADD COLUMN IF NOT EXISTS service_id TEXT,
    ADD COLUMN IF NOT EXISTS user_phone TEXT,
    ADD COLUMN IF NOT EXISTS provider_phone TEXT;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_services_point_code ON services (point_code);
CREATE INDEX IF NOT EXISTS idx_services_category_id ON services (category_id);
CREATE INDEX IF NOT EXISTS idx_services_subcategory_id ON services (subcategory_id);
CREATE INDEX IF NOT EXISTS idx_master_services_master_phone ON master_services (master_phone);
CREATE INDEX IF NOT EXISTS idx_master_services_service_id ON master_services (service_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE bagsies
DROP COLUMN IF EXISTS service_id,
    DROP COLUMN IF EXISTS user_phone,
    DROP COLUMN IF EXISTS provider_phone;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN IF EXISTS point_code;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS master_services;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS services;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS service_subcategories;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS service_categories;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE points
DROP COLUMN IF EXISTS category_id,
    DROP COLUMN IF EXISTS active,
    DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd

-- +goose StatementBegin
DROP TABLE IF EXISTS point_categories;
-- +goose StatementEnd
