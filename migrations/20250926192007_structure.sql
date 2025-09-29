-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS networks
(
    code        TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    description TEXT
);

ALTER TABLE points
    ADD COLUMN IF NOT EXISTS network_code TEXT,
    ADD COLUMN IF NOT EXISTS address      TEXT,
    ADD COLUMN IF NOT EXISTS city         TEXT;

ALTER TABLE bagsies
DROP CONSTRAINT IF EXISTS bagsies_pkey,
DROP COLUMN IF EXISTS id,
    ADD COLUMN id text PRIMARY KEY,
    ADD COLUMN IF NOT EXISTS start_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS end_at   TIMESTAMP;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS networks;
ALTER TABLE bagsies
DROP COLUMN IF EXISTS start_at,
    DROP COLUMN IF EXISTS end_at,
    DROP COLUMN IF EXISTS id,
    ADD COLUMN id SERIAL PRIMARY KEY;

ALTER TABLE points
DROP COLUMN IF EXISTS network_code,
    DROP COLUMN IF EXISTS address,
    DROP COLUMN IF EXISTS city;
-- +goose StatementEnd
