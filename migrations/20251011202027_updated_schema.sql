-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS point_categories (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE points
    ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES point_categories(id),
    ADD COLUMN IF NOT EXISTS active      BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS deleted_at  TIMESTAMPTZ;

CREATE TABLE IF NOT EXISTS service_categories (
    id          BIGSERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS service_subcategories (
    id          BIGSERIAL PRIMARY KEY,
    category_id BIGINT REFERENCES service_categories(id) ON DELETE CASCADE,
    name        TEXT NOT NULL,
    description TEXT,
    created_at  TIMESTAMPTZ DEFAULT now(),
    updated_at  TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS services (
    id               BIGSERIAL PRIMARY KEY,
    point_code       TEXT REFERENCES points(code) ON DELETE CASCADE,
    name             TEXT NOT NULL,
    description      TEXT,
    category_id      BIGINT REFERENCES service_categories(id),
    subcategory_id   BIGINT REFERENCES service_subcategories(id),
    duration_minutes TEXT,
    active           BOOLEAN DEFAULT TRUE,
    created_at       TIMESTAMPTZ DEFAULT now(),
    updated_at       TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS master_services (
    id               BIGSERIAL PRIMARY KEY,
    master_phone     TEXT REFERENCES users(phone) ON DELETE CASCADE,
    service_id       BIGINT REFERENCES services(id) ON DELETE CASCADE,
    price            TEXT,
    duration_minutes TEXT,
    active           BOOLEAN DEFAULT TRUE,
    created_at       TIMESTAMPTZ DEFAULT now(),
    updated_at       TIMESTAMPTZ DEFAULT now()
);

ALTER TABLE users
    ADD CONSTRAINT fk_users_point FOREIGN KEY (point_code) REFERENCES points(code) ON DELETE SET NULL;

ALTER TABLE bagsies
    ADD COLUMN IF NOT EXISTS service_id BIGINT REFERENCES services(id),
    ADD CONSTRAINT fk_bagsies_user FOREIGN KEY (user_phone) REFERENCES users(phone) ON DELETE CASCADE,
    ADD CONSTRAINT fk_bagsies_provider FOREIGN KEY (provider_phone) REFERENCES users(phone) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_services_point_code ON services (point_code);
CREATE INDEX IF NOT EXISTS idx_services_category_id ON services (category_id);
CREATE INDEX IF NOT EXISTS idx_services_subcategory_id ON services (subcategory_id);
CREATE INDEX IF NOT EXISTS idx_master_services_master_phone ON master_services (master_phone);
CREATE INDEX IF NOT EXISTS idx_master_services_service_id ON master_services (service_id);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE bagsies
DROP CONSTRAINT IF EXISTS fk_bagsies_user,
    DROP CONSTRAINT IF EXISTS fk_bagsies_provider,
    DROP COLUMN IF EXISTS service_id;

ALTER TABLE users
DROP CONSTRAINT IF EXISTS fk_users_point;

DROP TABLE IF EXISTS master_services;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS service_subcategories;
DROP TABLE IF EXISTS service_categories;

ALTER TABLE points
DROP COLUMN IF EXISTS category_id,
    DROP COLUMN IF EXISTS active,
    DROP COLUMN IF EXISTS deleted_at;

DROP TABLE IF EXISTS point_categories;
-- +goose StatementEnd
