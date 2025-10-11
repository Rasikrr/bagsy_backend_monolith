-- +goose Up
-- +goose StatementBegin
ALTER TABLE networks
    ADD COLUMN updated_at TIMESTAMPTZ,
    ADD COLUMN deleted_at TIMESTAMPTZ,
    ADD COLUMN created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_by TEXT;

ALTER TABLE points DROP COLUMN address;
ALTER TABLE points RENAME COLUMN coordinates TO address;
ALTER TABLE points
ALTER COLUMN created_at TYPE TIMESTAMPTZ,
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
    ADD CONSTRAINT point_code_network_code_unique UNIQUE (code, network_code);

ALTER TABLE services DROP COLUMN id;
ALTER TABLE services ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();
ALTER TABLE services
ALTER COLUMN duration_minutes TYPE INTEGER USING duration_minutes::integer,
    ADD CONSTRAINT services_id_cat_subcat_unique UNIQUE (id, category_id, subcategory_id);

ALTER TABLE master_services
DROP COLUMN id,
    DROP COLUMN duration_minutes;
ALTER TABLE master_services
    ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();

ALTER TABLE bagsies
DROP COLUMN first_name,
    DROP COLUMN last_name,
    DROP COLUMN service,
    DROP COLUMN id;

ALTER TABLE bagsies RENAME COLUMN provider_phone TO master_phone;

ALTER TABLE bagsies
ALTER COLUMN created_at TYPE TIMESTAMPTZ,
    ALTER COLUMN updated_at TYPE TIMESTAMPTZ,
    ALTER COLUMN start_at TYPE TIMESTAMPTZ,
    ALTER COLUMN end_at TYPE TIMESTAMPTZ,
    ADD COLUMN id UUID PRIMARY KEY DEFAULT gen_random_uuid();

ALTER TABLE users DROP COLUMN phone;
ALTER TABLE users ADD COLUMN phone TEXT PRIMARY KEY;

UPDATE networks SET updated_by = 'system' WHERE updated_by IS NULL;
ALTER TABLE networks
    ALTER COLUMN updated_by SET DEFAULT 'system',
ALTER COLUMN updated_by SET NOT NULL;

UPDATE points SET updated_by = 'system' WHERE updated_by IS NULL;
ALTER TABLE points
    ALTER COLUMN updated_by SET DEFAULT 'system',
ALTER COLUMN updated_by SET NOT NULL;

UPDATE bagsies SET updated_by = 'system' WHERE updated_by IS NULL;
ALTER TABLE bagsies
    ALTER COLUMN updated_by SET DEFAULT 'system',
ALTER COLUMN updated_by SET NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE networks
DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS deleted_at,
    DROP COLUMN IF EXISTS created_at,
    DROP COLUMN IF EXISTS updated_by;

ALTER TABLE points DROP CONSTRAINT IF EXISTS point_code_network_code_unique;
ALTER TABLE points
ALTER COLUMN created_at TYPE TIMESTAMP,
    ALTER COLUMN updated_at TYPE TIMESTAMP;
ALTER TABLE points RENAME COLUMN address TO coordinates;
ALTER TABLE points ADD COLUMN address TEXT;

ALTER TABLE services DROP CONSTRAINT IF EXISTS services_id_cat_subcat_unique;
ALTER TABLE services DROP COLUMN IF EXISTS id;
ALTER TABLE services ADD COLUMN id BIGSERIAL PRIMARY KEY;
ALTER TABLE services ALTER COLUMN duration_minutes TYPE TEXT;

ALTER TABLE master_services DROP COLUMN IF EXISTS id;
ALTER TABLE master_services ADD COLUMN id BIGSERIAL PRIMARY KEY;
ALTER TABLE master_services ADD COLUMN duration_minutes TEXT;

ALTER TABLE bagsies DROP COLUMN IF EXISTS id;
ALTER TABLE bagsies ADD COLUMN id TEXT PRIMARY KEY;
ALTER TABLE bagsies
ALTER COLUMN created_at TYPE TIMESTAMP,
    ALTER COLUMN updated_at TYPE TIMESTAMP,
    ALTER COLUMN start_at TYPE TIMESTAMP,
    ALTER COLUMN end_at TYPE TIMESTAMP;
ALTER TABLE bagsies RENAME COLUMN master_phone TO provider_phone;
ALTER TABLE bagsies
    ADD COLUMN first_name TEXT,
    ADD COLUMN last_name TEXT,
    ADD COLUMN service TEXT;

ALTER TABLE users DROP COLUMN IF EXISTS phone;
ALTER TABLE users ADD COLUMN phone TEXT UNIQUE;
-- +goose StatementEnd
