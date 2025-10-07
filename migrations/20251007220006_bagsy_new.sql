-- +goose Up
-- +goose StatementBegin
ALTER TABLE bagsies
    -- заменяем старые поля phone -> user_phone, если еще не переименовано
    RENAME COLUMN phone TO user_phone;

-- добавляем недостающие поля, если их нет
ALTER TABLE bagsies
    ADD COLUMN IF NOT EXISTS provider_phone TEXT,
    ADD COLUMN IF NOT EXISTS first_name     TEXT,
    ADD COLUMN IF NOT EXISTS last_name      TEXT,
    ADD COLUMN IF NOT EXISTS description    TEXT,
    ADD COLUMN IF NOT EXISTS service        TEXT;

-- приведение типов start_at / end_at, если они вдруг без NOT NULL
ALTER TABLE bagsies
ALTER COLUMN start_at TYPE TIMESTAMP USING start_at::timestamp,
    ALTER COLUMN end_at   TYPE TIMESTAMP USING end_at::timestamp;

-- индексы для ускорения поиска
CREATE INDEX IF NOT EXISTS idx_bagsies_user_phone     ON bagsies (user_phone);
CREATE INDEX IF NOT EXISTS idx_bagsies_provider_phone ON bagsies (provider_phone);
CREATE INDEX IF NOT EXISTS idx_bagsies_start_at       ON bagsies (start_at);
CREATE INDEX IF NOT EXISTS idx_bagsies_end_at         ON bagsies (end_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- откатываем всё, кроме rename (его обратный шаг может поломать совместимость)
ALTER TABLE bagsies
DROP COLUMN IF EXISTS provider_phone,
    DROP COLUMN IF EXISTS first_name,
    DROP COLUMN IF EXISTS last_name,
    DROP COLUMN IF EXISTS description,
    DROP COLUMN IF EXISTS service;

DROP INDEX IF EXISTS idx_bagsies_user_phone;
DROP INDEX IF EXISTS idx_bagsies_provider_phone;
DROP INDEX IF EXISTS idx_bagsies_start_at;
DROP INDEX IF EXISTS idx_bagsies_end_at;
-- +goose StatementEnd
